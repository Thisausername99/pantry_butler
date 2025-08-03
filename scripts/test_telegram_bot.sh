#!/bin/bash

# Telegram Bot Test Script
# This script helps test your Telegram bot configuration

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required tools are installed
check_dependencies() {
    print_status "Checking dependencies..."
    
    if ! command -v curl &> /dev/null; then
        print_error "curl is not installed. Please install curl first."
        exit 1
    fi
    
    if ! command -v jq &> /dev/null; then
        print_warning "jq is not installed. JSON parsing will be limited."
    fi
    
    print_success "Dependencies check completed"
}

# Test bot token
test_bot_token() {
    local token=$1
    
    if [ -z "$token" ]; then
        print_error "Bot token is required. Please provide it as an argument."
        echo "Usage: $0 <bot_token> [chat_id]"
        exit 1
    fi
    
    print_status "Testing bot token..."
    
    local response=$(curl -s "https://api.telegram.org/bot$token/getMe")
    
    if command -v jq &> /dev/null; then
        local ok=$(echo "$response" | jq -r '.ok')
        if [ "$ok" = "true" ]; then
            local bot_name=$(echo "$response" | jq -r '.result.first_name')
            local bot_username=$(echo "$response" | jq -r '.result.username')
            print_success "Bot token is valid!"
            print_status "Bot name: $bot_name"
            print_status "Bot username: @$bot_username"
        else
            print_error "Invalid bot token"
            print_error "Response: $response"
            exit 1
        fi
    else
        if echo "$response" | grep -q '"ok":true'; then
            print_success "Bot token is valid!"
        else
            print_error "Invalid bot token"
            print_error "Response: $response"
            exit 1
        fi
    fi
}

# Test chat ID
test_chat_id() {
    local token=$1
    local chat_id=$2
    
    if [ -z "$chat_id" ]; then
        print_warning "Chat ID not provided. Skipping chat ID test."
        return
    fi
    
    print_status "Testing chat ID..."
    
    local response=$(curl -s "https://api.telegram.org/bot$token/getChat" \
        -d "chat_id=$chat_id")
    
    if command -v jq &> /dev/null; then
        local ok=$(echo "$response" | jq -r '.ok')
        if [ "$ok" = "true" ]; then
            local chat_title=$(echo "$response" | jq -r '.result.title // .result.first_name // "Unknown"')
            local chat_type=$(echo "$response" | jq -r '.result.type')
            print_success "Chat ID is valid!"
            print_status "Chat title: $chat_title"
            print_status "Chat type: $chat_type"
        else
            local error_code=$(echo "$response" | jq -r '.error_code')
            local description=$(echo "$response" | jq -r '.description')
            print_error "Invalid chat ID"
            print_error "Error code: $error_code"
            print_error "Description: $description"
        fi
    else
        if echo "$response" | grep -q '"ok":true'; then
            print_success "Chat ID is valid!"
        else
            print_error "Invalid chat ID"
            print_error "Response: $response"
        fi
    fi
}

# Send test message
send_test_message() {
    local token=$1
    local chat_id=$2
    
    if [ -z "$chat_id" ]; then
        print_warning "Chat ID not provided. Skipping test message."
        return
    fi
    
    print_status "Sending test message..."
    
    local message="🧪 Test Message from Pantry Butler Bot

This is a test message to verify your Telegram bot integration is working correctly.

Timestamp: $(date -u +"%Y-%m-%d %H:%M:%S UTC")
Repository: pantry_butler
Test Type: Manual verification

✅ If you receive this message, your bot is configured correctly!"

    local response=$(curl -s "https://api.telegram.org/bot$token/sendMessage" \
        -d "chat_id=$chat_id" \
        -d "text=$message")
    
    if command -v jq &> /dev/null; then
        local ok=$(echo "$response" | jq -r '.ok')
        if [ "$ok" = "true" ]; then
            local message_id=$(echo "$response" | jq -r '.result.message_id')
            print_success "Test message sent successfully!"
            print_status "Message ID: $message_id"
        else
            local error_code=$(echo "$response" | jq -r '.error_code')
            local description=$(echo "$response" | jq -r '.description')
            print_error "Failed to send test message"
            print_error "Error code: $error_code"
            print_error "Description: $description"
        fi
    else
        if echo "$response" | grep -q '"ok":true'; then
            print_success "Test message sent successfully!"
        else
            print_error "Failed to send test message"
            print_error "Response: $response"
        fi
    fi
}

# Get updates (useful for debugging)
get_updates() {
    local token=$1
    
    print_status "Getting recent updates..."
    
    local response=$(curl -s "https://api.telegram.org/bot$token/getUpdates")
    
    if command -v jq &> /dev/null; then
        local ok=$(echo "$response" | jq -r '.ok')
        if [ "$ok" = "true" ]; then
            local update_count=$(echo "$response" | jq '.result | length')
            print_success "Found $update_count recent updates"
            
            if [ "$update_count" -gt 0 ]; then
                echo "$response" | jq '.result[] | {update_id: .update_id, message: .message.text // "No text", chat_id: .message.chat.id // "N/A"}'
            fi
        else
            print_error "Failed to get updates"
            print_error "Response: $response"
        fi
    else
        print_warning "jq not available. Raw response:"
        echo "$response"
    fi
}

# Main function
main() {
    echo "🧪 Telegram Bot Test Script"
    echo "=========================="
    echo ""
    
    local bot_token=$1
    local chat_id=$2
    
    check_dependencies
    echo ""
    
    test_bot_token "$bot_token"
    echo ""
    
    test_chat_id "$bot_token" "$chat_id"
    echo ""
    
    send_test_message "$bot_token" "$chat_id"
    echo ""
    
    if [ "$3" = "--debug" ]; then
        get_updates "$bot_token"
        echo ""
    fi
    
    print_success "Test completed!"
    echo ""
    print_status "Next steps:"
    echo "1. Add your bot token to GitHub secrets as TELEGRAM_BOT_TOKEN"
    echo "2. Add your chat ID to GitHub secrets as TELEGRAM_CHAT_ID"
    echo "3. Push a change to trigger the workflow"
    echo "4. Check your Telegram for notifications"
}

# Check if script is being sourced or executed
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi 