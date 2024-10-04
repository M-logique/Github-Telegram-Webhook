package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func formatGitHubWebhook(r *http.Request, blacklistedActions string, blacklistedevents string) (string, error) {
	var webhook GitHubWebhook
	var targetedAction string = webhook.Action

	
	if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
		return "", err
	}
	event := r.Header.Get("X-GitHub-Event")

	// Action can be empty, so we need to use this if
	if webhook.Action != "" {
		// Actions can be duplicated; for more customization, it's better to add
		// the event at the end of the action
		targetedAction = fmt.Sprintf("%s_%s", event, webhook.Action)
	}

	if isBlacklisted(targetedAction, blacklistedActions) {
		return "", fmt.Errorf("action is blacklisted")
	}

	if isBlacklisted(event, blacklistedevents) {
		return "", fmt.Errorf("event is blacklisted")
	}
	title := getEventTitle(webhook)
	details := formatEventDetails(event, webhook)
	return fmt.Sprintf("%s\n%s", title, details), nil
}

func getEventTitle(webhook GitHubWebhook) string {
	repoName := webhook.Repository.FullName
	sender := webhook.Sender.Login
	action := webhook.Action
	return fmt.Sprintf("<b>%s</b> | <i>%s %s</i>", repoName, sender, action)
}

func formatEventDetails(event string, webhook GitHubWebhook) string {
	switch event {
	case "create":
		return fmt.Sprintf("✨ → <a href=\"%s/tree/%s\">%s: %s</a> created by %s",
			webhook.Repository.HTMLURL, webhook.Ref, webhook.RefType, webhook.Ref, webhook.Sender.Login)
	case "delete":
		return fmt.Sprintf("🗑️ → %s: <code>%s</code> deleted by %s", webhook.RefType, webhook.Ref, webhook.Sender.Login)
	case "push":
		return fmt.Sprintf("🚀 → %s pushed to <a href=\"%s/commits/%s\">%s</a> by %s",
			webhook.Ref, webhook.Repository.HTMLURL, webhook.Ref, webhook.Ref, webhook.Sender.Login)
	case "pull_request":
		return fmt.Sprintf("🔄 → Pull request %s by %s on <a href=\"%s/pull/%s\">%s</a>",
			webhook.Action, webhook.Sender.Login, webhook.Repository.HTMLURL, webhook.Ref, webhook.Ref)
	case "fork":
		return fmt.Sprintf("🍴 → %s forked the repository: <a href=\"%s\">%s</a>",
			webhook.Sender.Login, webhook.Repository.HTMLURL, webhook.Repository.FullName)
	case "release":
		return fmt.Sprintf("🎉 → Release %s by %s on <a href=\"%s/releases/tag/%s\">%s</a>",
			webhook.Action, webhook.Sender.Login, webhook.Repository.HTMLURL, webhook.Ref, webhook.Ref)
	case "workflow_run":
		return fmt.Sprintf("🏃 → Workflow run: %s by %s", webhook.Action, webhook.Sender.Login)
	case "workflow_job":
		return fmt.Sprintf("🔧 → Workflow job: %s by %s", webhook.Action, webhook.Sender.Login)
	case "check_run":
		return fmt.Sprintf("✅ → Check run: %s by %s", webhook.Action, webhook.Sender.Login)
	case "check_suite":
		return fmt.Sprintf("🔍 → Check suite: %s by %s", webhook.Action, webhook.Sender.Login)
	case "status":
		return fmt.Sprintf("📊 → Status event: %s by %s", webhook.Action, webhook.Sender.Login)
	case "issue":
		return fmt.Sprintf("🐛 → Issue %s: <a href=\"%s/issues/%s\">#%s</a> opened by %s",
			webhook.Action, webhook.Repository.HTMLURL, webhook.Ref, webhook.Ref, webhook.Sender.Login)
	case "branch_protection_rule":
		return fmt.Sprintf("🛡️ → Branch protection rule updated by %s", webhook.Sender.Login)
	case "code_scanning_alert":
		return fmt.Sprintf("🛠️ → Code scanning alert %s by %s", webhook.Action, webhook.Sender.Login)
	case "commit_comment":
		return fmt.Sprintf("💬 → Commit comment by %s on <a href=\"%s/commit/%s\">commit</a>", webhook.Sender.Login, webhook.Repository.HTMLURL, webhook.Ref)
	case "deployment":
		return fmt.Sprintf("📦 → Deployment created by %s", webhook.Sender.Login)
	case "deployment_status":
		return fmt.Sprintf("🔄 → Deployment status: %s by %s", webhook.Action, webhook.Sender.Login)
	case "discussion":
		return fmt.Sprintf("💬 → Discussion %s by %s", webhook.Action, webhook.Sender.Login)
	case "discussion_comment":
		return fmt.Sprintf("🗣️ → Discussion comment %s by %s", webhook.Action, webhook.Sender.Login)
	case "gollum":
		return fmt.Sprintf("📚 → Wiki page updated by %s", webhook.Sender.Login)
	case "issues":
		return fmt.Sprintf("🐞 → Issue %s by %s on <a href=\"%s/issues/%s\">#%s</a>", webhook.Action, webhook.Sender.Login, webhook.Repository.HTMLURL, webhook.Ref, webhook.Ref)
	case "issue_comment":
		return fmt.Sprintf("💬 → Issue comment %s by %s", webhook.Action, webhook.Sender.Login)
	case "label":
		return fmt.Sprintf("🏷️ → Label %s by %s", webhook.Action, webhook.Sender.Login)
	case "member":
		return fmt.Sprintf("👥 → Member %s by %s", webhook.Action, webhook.Sender.Login)
	case "membership":
		return fmt.Sprintf("👤 → Membership %s by %s", webhook.Action, webhook.Sender.Login)
	case "merge_group":
		return fmt.Sprintf("🔗 → Merge group %s by %s", webhook.Action, webhook.Sender.Login)
	case "milestone":
		return fmt.Sprintf("🎯 → Milestone %s by %s", webhook.Action, webhook.Sender.Login)
	case "organization":
		return fmt.Sprintf("🏢 → Organization event %s by %s", webhook.Action, webhook.Sender.Login)
	case "org_block":
		return fmt.Sprintf("⛔ → Organization block %s by %s", webhook.Action, webhook.Sender.Login)
	case "package":
		return fmt.Sprintf("📦 → Package event %s by %s", webhook.Action, webhook.Sender.Login)
	case "page_build":
		return fmt.Sprintf("🌐 → Page build %s by %s", webhook.Action, webhook.Sender.Login)
	case "ping":
		return fmt.Sprintf("📶 → Ping event triggered by %s", webhook.Sender.Login)
	case "project":
		return fmt.Sprintf("📁 → Project %s by %s", webhook.Action, webhook.Sender.Login)
	case "project_card":
		return fmt.Sprintf("📝 → Project card %s by %s", webhook.Action, webhook.Sender.Login)
	case "project_column":
		return fmt.Sprintf("📊 → Project column %s by %s", webhook.Action, webhook.Sender.Login)
	case "public":
		return fmt.Sprintf("🌍 → Repository %s made public by %s", webhook.Action, webhook.Sender.Login)
	case "pull_request_review":
		return fmt.Sprintf("📝 → Pull request review %s by %s", webhook.Action, webhook.Sender.Login)
	case "pull_request_review_comment":
		return fmt.Sprintf("💬 → Pull request review comment %s by %s", webhook.Action, webhook.Sender.Login)
	case "registry_package":
		return fmt.Sprintf("📦 → Registry package %s by %s", webhook.Action, webhook.Sender.Login)
	case "repository":
		return fmt.Sprintf("📂 → Repository %s by %s", webhook.Action, webhook.Sender.Login)
	case "repository_dispatch":
		return fmt.Sprintf("📤 → Repository dispatch %s by %s", webhook.Action, webhook.Sender.Login)
	case "repository_import":
		return fmt.Sprintf("📥 → Repository import %s by %s", webhook.Action, webhook.Sender.Login)
	case "repository_vulnerability_alert":
		return fmt.Sprintf("⚠️ → Vulnerability alert %s by %s", webhook.Action, webhook.Sender.Login)
	case "secret_scanning_alert":
		return fmt.Sprintf("🔐 → Secret scanning alert %s by %s", webhook.Action, webhook.Sender.Login)
	case "secret_scanning_alert_location":
		return fmt.Sprintf("🔍 → Secret scanning alert location %s by %s", webhook.Action, webhook.Sender.Login)
	case "security_advisory":
		return fmt.Sprintf("🔒 → Security advisory %s by %s", webhook.Action, webhook.Sender.Login)
	case "sponsorship":
		return fmt.Sprintf("💰 → Sponsorship %s by %s", webhook.Action, webhook.Sender.Login)
	case "star":
		return fmt.Sprintf("⭐ → Star event %s by %s", webhook.Action, webhook.Sender.Login)
	case "team":
		return fmt.Sprintf("👥 → Team event %s by %s", webhook.Action, webhook.Sender.Login)
	case "team_add":
		return fmt.Sprintf("➕ → Team add event %s by %s", webhook.Action, webhook.Sender.Login)
	case "watch":
		return fmt.Sprintf("👁️ → Watch event %s by %s", webhook.Action, webhook.Sender.Login)
	default:
		return fmt.Sprintf("❓ → Unknown event: %s", event)
	}
}