// Package action provides MCP tools related to Forgejo Actions.
//
// It currently implements the `list_action_tasks`, `list_action_run_jobs`,
// `get_action_job_logs` and `dispatch_workflow` tools, which allow listing
// Forgejo Actions execution tasks, inspecting the jobs of a workflow run,
// downloading job logs and triggering workflows via the workflow_dispatch
// event, extending functionality beyond the official Forgejo SDK.
package action
