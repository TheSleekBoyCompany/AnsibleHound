package opengraph

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"slices"
	"sort"
	"sync"

	"github.com/Ramoreik/awxlimit/pkg/awxlimit"
	"github.com/Ramoreik/gopengraph"
	"github.com/Ramoreik/gopengraph/edge"
	"github.com/Ramoreik/gopengraph/node"
	"github.com/charmbracelet/log"
)

type credentialWork struct {
	ordinal               int
	credentialID          string
	credentialType        string
	machineCredentialType string
	hasEnvInjectors       bool
	hasExtraVarInjectors  bool
	relations             []credentialRelation
}

type credentialRelation struct {
	kind                  string
	startID               string
	canControlPlaybook    bool
	canControlJobTemplate bool
	canControlInventory   bool
	inventory             *inventoryWork
}

type inventoryWork struct {
	limit     string
	inventory awxlimit.Inventory
	targets   []inventoryTarget
}

type inventoryTarget struct {
	id   string
	name string
}

type edgeCandidate struct {
	ordinal int
	kind    string
	startID string
	endID   string
}

type credentialResult struct {
	candidates []edgeCandidate
	err        error
}

// PostProcessingCredentials derives credential attack paths concurrently and commits them serially.
//
// OpenGraph is not safe for concurrent access. This function first snapshots every graph value
// required for credential derivation, then runs graph-free workers, and finally adds the resulting
// edges after all workers have completed.
func PostProcessingCredentials(ctx context.Context, graph *gopengraph.OpenGraph, maxWorkers int) error {
	if ctx == nil {
		return errors.New("post-process credentials: nil context")
	}

	log.Info("Handling post processing edges for Credentials.")

	work, err := buildCredentialWork(graph)
	if err != nil {
		return fmt.Errorf("snapshot credential post-processing work: %w", err)
	}

	candidates, err := deriveCredentialCandidates(ctx, work, maxWorkers)
	if err != nil {
		return fmt.Errorf("derive credential post-processing edges: %w", err)
	}

	if err := commitCredentialCandidates(graph, candidates); err != nil {
		return fmt.Errorf("commit credential post-processing edges: %w", err)
	}

	return nil
}

func buildCredentialWork(graph *gopengraph.OpenGraph) ([]credentialWork, error) {
	credentialNodes := graph.GetNodesByKind(CREDENTIAL_NODE)
	sort.Slice(credentialNodes, func(i, j int) bool {
		return credentialNodes[i].GetID() < credentialNodes[j].GetID()
	})

	work := make([]credentialWork, 0, len(credentialNodes))
	for ordinal, credentialNode := range credentialNodes {
		credentialType, credentialTypeNode, err := credentialTypeFor(graph, credentialNode)
		if err != nil {
			return nil, fmt.Errorf("credential %q: %w", credentialNode.GetID(), err)
		}

		item := credentialWork{
			ordinal:        ordinal,
			credentialID:   credentialNode.GetID(),
			credentialType: credentialType,
		}

		if credentialType == MACHINE_CREDENTIAL_TYPE {
			item.machineCredentialType, err = requiredString(credentialNode, "machine_credential_type")
			if err != nil {
				return nil, fmt.Errorf("credential %q: %w", credentialNode.GetID(), err)
			}
		} else {
			item.hasEnvInjectors, err = optionalBool(credentialTypeNode, "injector_env")
			if err != nil {
				return nil, fmt.Errorf("credential type %q: %w", credentialTypeNode.GetID(), err)
			}
			item.hasExtraVarInjectors, err = optionalBool(credentialTypeNode, "injector_extra_vars")
			if err != nil {
				return nil, fmt.Errorf("credential type %q: %w", credentialTypeNode.GetID(), err)
			}
		}

		for _, relationEdge := range graph.GetEdgesToNode(credentialNode.GetID()) {
			relation, err := snapshotCredentialRelation(graph, credentialType, relationEdge)
			if err != nil {
				return nil, fmt.Errorf("credential %q: %w", credentialNode.GetID(), err)
			}
			item.relations = append(item.relations, relation)
		}

		work = append(work, item)
	}

	return work, nil
}

func credentialTypeFor(graph *gopengraph.OpenGraph, credentialNode *node.Node) (string, *node.Node, error) {
	for _, relationship := range graph.GetEdgesFromNode(credentialNode.GetID()) {
		if relationship.GetKind() != USES_TYPE_EDGE {
			continue
		}

		credentialTypeNode := graph.GetNode(relationship.GetEndNodeID())
		if credentialTypeNode == nil {
			return "", nil, fmt.Errorf("credential type node %q does not exist", relationship.GetEndNodeID())
		}

		credentialType, err := requiredString(credentialTypeNode, "name")
		if err != nil {
			return "", nil, err
		}
		return credentialType, credentialTypeNode, nil
	}

	return "", nil, errors.New("credential type relationship does not exist")
}

func snapshotCredentialRelation(graph *gopengraph.OpenGraph, credentialType string, relationship *edge.Edge) (credentialRelation, error) {
	startNode := graph.GetNode(relationship.GetStartNodeID())
	if startNode == nil {
		return credentialRelation{}, fmt.Errorf("relationship source node %q does not exist", relationship.GetStartNodeID())
	}

	relation := credentialRelation{
		kind:    relationship.GetKind(),
		startID: relationship.GetStartNodeID(),
	}

	if relation.kind == USE_ROLE_EDGE || relation.kind == ADMIN_ROLE_EDGE {
		controls, err := snapshotControls(graph, startNode)
		if err != nil {
			return credentialRelation{}, err
		}
		relation.canControlPlaybook = controls.canControlPlaybook
		relation.canControlJobTemplate = controls.canControlJobTemplate
		relation.canControlInventory = controls.canControlInventory
	}

	if credentialType == MACHINE_CREDENTIAL_TYPE && relation.kind == USES_EDGE && startNode.HasKind(JOB_TEMPLATE_NODE) {
		inventory, err := snapshotJobTemplateInventory(graph, startNode)
		if err != nil {
			return credentialRelation{}, err
		}
		relation.inventory = inventory
	}

	return relation, nil
}

type identityControls struct {
	canControlPlaybook    bool
	canControlJobTemplate bool
	canControlInventory   bool
}

func snapshotControls(graph *gopengraph.OpenGraph, identity *node.Node) (identityControls, error) {
	var controls identityControls
	for _, relationship := range graph.GetEdgesFromNode(identity.GetID()) {
		endNode := graph.GetNode(relationship.GetEndNodeID())
		if endNode == nil {
			return identityControls{}, fmt.Errorf("identity %q references missing node %q", identity.GetID(), relationship.GetEndNodeID())
		}

		if endNode.HasKind(PROJECT_NODE) && relationship.GetKind() == ADMIN_ROLE_EDGE {
			controls.canControlPlaybook = true
		}
		if endNode.HasKind(ORGANIZATION_NODE) && (relationship.GetKind() == ADMIN_ROLE_EDGE || relationship.GetKind() == "ATProjectAdmin") {
			controls.canControlPlaybook = true
			controls.canControlJobTemplate = true
		}
		if endNode.HasKind(JOB_TEMPLATE_NODE) && relationship.GetKind() == ADMIN_ROLE_EDGE {
			controls.canControlJobTemplate = true
		}
		if endNode.HasKind(INVENTORY_NODE) && relationship.GetKind() == ADMIN_ROLE_EDGE {
			controls.canControlInventory = true
		}
		if endNode.HasKind(ORGANIZATION_NODE) && (relationship.GetKind() == INVENTORY_ADMIN_ROLE_EDGE || relationship.GetKind() == ADMIN_ROLE_EDGE) {
			controls.canControlInventory = true
		}
	}
	return controls, nil
}

func snapshotJobTemplateInventory(graph *gopengraph.OpenGraph, jobTemplate *node.Node) (*inventoryWork, error) {
	limit, err := optionalString(jobTemplate, "limit")
	if err != nil {
		return nil, fmt.Errorf("job template %q: %w", jobTemplate.GetID(), err)
	}

	for _, relationship := range graph.GetEdgesFromNode(jobTemplate.GetID()) {
		if relationship.GetKind() != USES_EDGE {
			continue
		}
		inventoryNode := graph.GetNode(relationship.GetEndNodeID())
		if inventoryNode == nil {
			return nil, fmt.Errorf("job template %q references missing node %q", jobTemplate.GetID(), relationship.GetEndNodeID())
		}
		if inventoryNode.HasKind(INVENTORY_NODE) {
			return snapshotInventory(graph, inventoryNode, limit)
		}
	}

	return nil, nil
}

func snapshotInventory(graph *gopengraph.OpenGraph, inventoryNode *node.Node, limit string) (*inventoryWork, error) {
	work := &inventoryWork{limit: limit}

	for _, relationship := range graph.GetEdgesFromNode(inventoryNode.GetID()) {
		if relationship.GetKind() != CONTAINS_EDGE {
			continue
		}

		member := graph.GetNode(relationship.GetEndNodeID())
		if member == nil {
			return nil, fmt.Errorf("inventory %q references missing node %q", inventoryNode.GetID(), relationship.GetEndNodeID())
		}

		switch {
		case member.HasKind(HOST_NODE):
			name, err := requiredString(member, "name")
			if err != nil {
				return nil, fmt.Errorf("host %q: %w", member.GetID(), err)
			}
			work.inventory.Hosts = append(work.inventory.Hosts, name)
			work.targets = append(work.targets, inventoryTarget{id: member.GetID(), name: name})
		case member.HasKind(GROUP_NODE):
			group, err := snapshotGroup(graph, member)
			if err != nil {
				return nil, err
			}
			work.inventory.Groups = append(work.inventory.Groups, group)
		}
	}

	return work, nil
}

func snapshotGroup(graph *gopengraph.OpenGraph, groupNode *node.Node) (awxlimit.Group, error) {
	name, err := requiredString(groupNode, "name")
	if err != nil {
		return awxlimit.Group{}, fmt.Errorf("group %q: %w", groupNode.GetID(), err)
	}

	group := awxlimit.Group{Name: name}
	for _, relationship := range graph.GetEdgesFromNode(groupNode.GetID()) {
		if relationship.GetKind() != CONTAINS_EDGE {
			continue
		}
		host := graph.GetNode(relationship.GetEndNodeID())
		if host == nil || !host.HasKind(HOST_NODE) {
			continue
		}
		hostName, err := requiredString(host, "name")
		if err != nil {
			return awxlimit.Group{}, fmt.Errorf("host %q: %w", host.GetID(), err)
		}
		group.Hosts = append(group.Hosts, hostName)
	}
	return group, nil
}

func deriveCredentialCandidates(ctx context.Context, work []credentialWork, maxWorkers int) ([]edgeCandidate, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(work) == 0 {
		return nil, nil
	}

	workerCount := postProcessingWorkerCount(maxWorkers, len(work))
	jobs := make(chan credentialWork)
	results := make(chan credentialResult, workerCount)
	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var workers sync.WaitGroup
	workers.Add(workerCount)
	for range workerCount {
		go func() {
			defer workers.Done()
			for {
				select {
				case <-workerCtx.Done():
					return
				case item, ok := <-jobs:
					if !ok {
						return
					}
					candidates, err := deriveCredentialCandidatesForWork(workerCtx, item)
					result := credentialResult{candidates: candidates, err: err}
					select {
					case results <- result:
					case <-workerCtx.Done():
					}
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for _, item := range work {
			select {
			case jobs <- item:
			case <-workerCtx.Done():
				return
			}
		}
	}()

	go func() {
		workers.Wait()
		close(results)
	}()

	var candidates []edgeCandidate
	var errs []error
	for result := range results {
		if result.err != nil {
			errs = append(errs, result.err)
			cancel()
			continue
		}
		candidates = append(candidates, result.candidates...)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	return candidates, nil
}

func postProcessingWorkerCount(maxWorkers, workItems int) int {
	workerCount := maxWorkers
	if workerCount == 0 {
		workerCount = runtime.GOMAXPROCS(0)
	}
	if workerCount < 1 {
		workerCount = 1
	}
	if workerCount > workItems {
		return workItems
	}
	return workerCount
}

func deriveCredentialCandidatesForWork(ctx context.Context, item credentialWork) ([]edgeCandidate, error) {
	var candidates []edgeCandidate
	for _, relation := range item.relations {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		switch item.credentialType {
		case SCM_CREDENTIAL_TYPE:
			if isCredentialRole(relation.kind) && relation.canControlPlaybook {
				candidates = appendCandidate(candidates, item, COMPROMISE_WITH_FAKE_SCM_SERVER_POST_PROCESSING_EDGE, relation.startID)
			}
		case SECRET_SERVER_CREDENTIAL_TYPE, HASHICORP_VAULT_CREDENTIAL_TYPE:
			if relation.kind == ADMIN_ROLE_EDGE {
				candidates = appendCandidate(candidates, item, COMPROMISE_WITH_REQUESTBIN_POST_PROCESSING_EDGE, relation.startID)
			}
		case MACHINE_CREDENTIAL_TYPE:
			if relation.inventory != nil {
				validFor, err := validForCandidates(item, *relation.inventory)
				if err != nil {
					return nil, err
				}
				candidates = append(candidates, validFor...)
			}
			if !isCredentialRole(relation.kind) {
				continue
			}
			if item.machineCredentialType == MACHINE_CREDENTIAL_SSH_SUBTYPE && (relation.canControlJobTemplate || relation.canControlPlaybook) {
				candidates = appendCandidate(candidates, item, SSH_HIJACK_AGENT_POST_PROCESSING_EDGE, relation.startID)
			}
			if item.machineCredentialType == MACHINE_CREDENTIAL_PASSWORD_SUBTYPE && relation.canControlInventory {
				candidates = appendCandidate(candidates, item, COMPROMISE_WITH_HONEYPOT_POST_PROCESSING_EDGE, relation.startID)
			}
			if relation.canControlInventory || relation.canControlPlaybook {
				candidates = appendCandidate(candidates, item, CAN_USE_IN_ADHOC_COMMANDS_POST_PROCESSING_EDGE, relation.startID)
			}
		default:
			if isCredentialRole(relation.kind) && (item.hasEnvInjectors || item.hasExtraVarInjectors) && relation.canControlPlaybook {
				candidates = appendCandidate(candidates, item, COMPROMISE_WITH_PLAYBOOK_POST_PROCESSING_EDGE, relation.startID)
			}
		}
	}
	return candidates, nil
}

func validForCandidates(item credentialWork, inventory inventoryWork) ([]edgeCandidate, error) {
	targets := inventory.targets
	if inventory.limit != "" {
		matched, err := awxlimit.MatchHosts(inventory.limit, inventory.inventory)
		if err != nil {
			return nil, fmt.Errorf("match inventory limit %q for credential %q: %w", inventory.limit, item.credentialID, err)
		}
		targets = make([]inventoryTarget, 0, len(inventory.targets))
		for _, target := range inventory.targets {
			if slices.Contains(matched, target.name) {
				targets = append(targets, target)
			}
		}
	}

	candidates := make([]edgeCandidate, 0, len(targets))
	for _, target := range targets {
		candidates = appendCandidate(candidates, item, VALID_FOR_POST_PROCESSING_EDGE, item.credentialID, target.id)
	}
	return candidates, nil
}

func appendCandidate(candidates []edgeCandidate, item credentialWork, kind, startID string, endIDs ...string) []edgeCandidate {
	endID := item.credentialID
	if len(endIDs) > 0 {
		endID = endIDs[0]
	}
	return append(candidates, edgeCandidate{ordinal: item.ordinal, kind: kind, startID: startID, endID: endID})
}

func isCredentialRole(kind string) bool {
	return kind == USE_ROLE_EDGE || kind == ADMIN_ROLE_EDGE
}

func commitCredentialCandidates(graph *gopengraph.OpenGraph, candidates []edgeCandidate) error {
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].ordinal != candidates[j].ordinal {
			return candidates[i].ordinal < candidates[j].ordinal
		}
		if candidates[i].kind != candidates[j].kind {
			return candidates[i].kind < candidates[j].kind
		}
		if candidates[i].startID != candidates[j].startID {
			return candidates[i].startID < candidates[j].startID
		}
		return candidates[i].endID < candidates[j].endID
	})

	seen := make(map[string]struct{}, len(candidates))
	for _, candidate := range candidates {
		key := candidate.kind + "\x00" + candidate.startID + "\x00" + candidate.endID
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		e, err := edge.NewEdge(candidate.startID, candidate.endID, candidate.kind, MATCH_BY_ID, MATCH_BY_ID, ANSIBLE_BASE, ANSIBLE_BASE, nil)
		if err != nil {
			return fmt.Errorf("create %q edge from %q to %q: %w", candidate.kind, candidate.startID, candidate.endID, err)
		}
		AddEdge(graph, e)
	}
	return nil
}

func requiredString(n *node.Node, key string) (string, error) {
	value, ok := n.GetProperty(key).(string)
	if !ok || value == "" {
		return "", fmt.Errorf("property %q must be a non-empty string", key)
	}
	return value, nil
}

func optionalString(n *node.Node, key string) (string, error) {
	value := n.GetProperty(key, "")
	stringValue, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("property %q must be a string, got %T", key, value)
	}
	return stringValue, nil
}

func optionalBool(n *node.Node, key string) (bool, error) {
	value := n.GetProperty(key, false)
	boolValue, ok := value.(bool)
	if !ok {
		return false, fmt.Errorf("property %q must be a bool, got %T", key, value)
	}
	return boolValue, nil
}
