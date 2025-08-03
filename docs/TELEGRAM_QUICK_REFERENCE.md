# Telegram Bot Quick Reference

## 🚀 Quick Setup

1. **Create Bot**: Message `@BotFather` → `/newbot` → Follow prompts
2. **Get Chat ID**: Message `@userinfobot` or use the test script
3. **Add Secrets**: Go to repo Settings → Secrets → Add:
   - `TELEGRAM_BOT_TOKEN`: Your bot token
   - `TELEGRAM_CHAT_ID`: Your chat ID

## 🧪 Testing

```bash
# Test your bot configuration
./scripts/test_telegram_bot.sh YOUR_BOT_TOKEN YOUR_CHAT_ID

# Test with debug info
./scripts/test_telegram_bot.sh YOUR_BOT_TOKEN YOUR_CHAT_ID --debug
```

## 📱 Notification Types

| Event | Trigger | Message |
|-------|---------|---------|
| Push | Code pushed to main/master | 🚀 New Code Push |
| PR | Pull request opened/updated | 🔄 Pull Request [action] |
| CI/CD | Workflow completed | 🔧 CI/CD Pipeline [status] |
| Release | Release published | 🎉 Release [action] |
| Failure | Any workflow fails | ❌ Workflow Failed |
| Deployment | Main branch deployment | ✅ Deployment Successful |

## 🔧 Configuration

### Workflow File
- **Location**: `.github/workflows/telegram-notifications.yml`
- **Triggers**: Push, PR, Workflow Run, Release
- **Branches**: main, master, develop

### Customization
```yaml
# Add/remove branches
on:
  push:
    branches: [ main, master, develop, feature/* ]

# Customize message format
message: |
  🚀 **Custom Notification**
  **Repository:** `${{ github.repository }}`
  **Branch:** `${{ github.ref_name }}`
```

## 🐛 Troubleshooting

### No Notifications
- ✅ Bot token correct?
- ✅ Chat ID correct?
- ✅ Started conversation with bot?
- ✅ Secrets added to GitHub?
- ✅ Workflow file in correct location?

### Common Errors
- **403 Forbidden**: Bot blocked or no permission
- **400 Bad Request**: Invalid chat ID
- **401 Unauthorized**: Invalid bot token

### Debug Steps
1. Check GitHub Actions logs
2. Run test script locally
3. Verify bot can send messages manually
4. Check workflow trigger conditions

## 📊 GitHub Secrets

| Secret Name | Description | Required |
|-------------|-------------|----------|
| `TELEGRAM_BOT_TOKEN` | Bot token from BotFather | ✅ |
| `TELEGRAM_CHAT_ID` | Your chat ID | ✅ |
| `TELEGRAM_CHAT_ID_DEV` | Dev team chat ID | ❌ |
| `TELEGRAM_CHAT_ID_ADMIN` | Admin chat ID | ❌ |

## 🔗 Useful Links

- [BotFather](https://t.me/botfather) - Create Telegram bots
- [User Info Bot](https://t.me/userinfobot) - Get your chat ID
- [Telegram Bot API](https://core.telegram.org/bots/api) - API documentation
- [GitHub Actions](https://docs.github.com/en/actions) - Workflow documentation

## 📝 Example Messages

### Push Notification
```
🚀 New Code Push

Repository: `username/pantry_butler`
Branch: `main`
Commit: `abc123def`
Author: `developer`
Message: `feat: add new recipe endpoint`

📊 View Commit
```

### CI/CD Success
```
🔧 CI/CD Pipeline success

Repository: `username/pantry_butler`
Workflow: `CI/CD Pipeline`
Branch: `main`
Commit: `abc123def`
Status: success
Duration: 2m 30s

📊 View Run
```

## 🛡️ Security

- Never commit bot tokens to code
- Use repository secrets for sensitive data
- Regularly rotate bot tokens
- Monitor bot usage for unusual activity
- Limit bot permissions to minimum required

## 📞 Support

- Check [setup guide](TELEGRAM_BOT_SETUP.md) for detailed instructions
- Review GitHub Actions logs for errors
- Test bot manually before troubleshooting
- Open issue in repository for specific problems 