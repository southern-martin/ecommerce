#!/usr/bin/env zsh
# Watch background agent progress in real-time (zsh compatible)
# Usage: ./scripts/watch-agents.sh [agent_id1] [agent_id2] ...
# If no args, watches all active agents

TASKS_DIR="/private/tmp/claude-501/-opt-openAi-ecommerce/tasks"
GREEN="\033[1;32m"
RED="\033[1;31m"
YELLOW="\033[1;33m"
CYAN="\033[1;36m"
MAGENTA="\033[1;35m"
DIM="\033[2m"
BOLD="\033[1m"
RESET="\033[0m"

typeset -a AGENT_IDS
typeset -a AGENT_FILES
typeset -a LAST_LINES
typeset -a LAST_ACTIVE
typeset -a AGENT_COLORS

ALL_COLORS=("$CYAN" "$YELLOW" "$MAGENTA" "$GREEN" "$RED")

# Collect agents
if [[ $# -gt 0 ]]; then
    for id in "$@"; do
        f="$TASKS_DIR/$id.output"
        if [[ -f "$f" ]]; then
            AGENT_IDS+=("$id")
            AGENT_FILES+=("$f")
        else
            echo "Warning: No output file for agent $id"
        fi
    done
else
    for f in "$TASKS_DIR"/*.output(N); do
        [[ -f "$f" ]] || continue
        id="${f:t:r}"
        AGENT_IDS+=("$id")
        AGENT_FILES+=("$f")
    done
fi

if [[ ${#AGENT_IDS[@]} -eq 0 ]]; then
    echo "No agent output files found in $TASKS_DIR"
    exit 1
fi

# Initialize tracking arrays
for i in {1..${#AGENT_IDS[@]}}; do
    LAST_LINES+=("$(wc -l < "${AGENT_FILES[$i]}" 2>/dev/null || echo 0)")
    LAST_ACTIVE+=("$(date +%s)")
    AGENT_COLORS+=("${ALL_COLORS[$(( (i - 1) % ${#ALL_COLORS[@]} + 1 ))]}")
done

echo -e "${BOLD}Monitoring ${#AGENT_IDS[@]} agent(s)${RESET}"
echo -e "${DIM}Press Ctrl+C to stop${RESET}"
echo ""
for i in {1..${#AGENT_IDS[@]}}; do
    echo -e "${AGENT_COLORS[$i]}[${AGENT_IDS[$i]}]${RESET} watching ${AGENT_FILES[$i]}"
done
echo "─────────────────────────────────────────────────────"

process_line() {
    local line="$1"
    local agent_id="$2"
    local color="$3"

    # File created
    if [[ "$line" == *"File created successfully"* ]]; then
        local fp="${line##*File created successfully at: }"
        fp="${fp%%\"*}"
        fp="${fp#/opt/openAi/ecommerce/}"
        [[ -n "$fp" ]] && echo -e "${color}[${agent_id}]${RESET} ${GREEN}✓ Created${RESET} $fp"
        return 0
    fi

    # File updated
    if [[ "$line" == *"has been updated successfully"* ]]; then
        local fp=""
        if [[ "$line" =~ '/opt/openAi/ecommerce/[^"]*' ]]; then
            fp="${MATCH#/opt/openAi/ecommerce/}"
        fi
        [[ -n "$fp" ]] && echo -e "${color}[${agent_id}]${RESET} ${GREEN}✎ Updated${RESET} $fp"
        return 0
    fi

    # Writing file
    if [[ "$line" == *'"name":"Write"'* && "$line" == *'"file_path"'* ]]; then
        local fp=""
        if [[ "$line" =~ '"file_path":"([^"]*)"' ]]; then
            fp="${match[1]#/opt/openAi/ecommerce/}"
        fi
        [[ -n "$fp" ]] && echo -e "${color}[${agent_id}]${RESET} ${DIM}⟳ Writing${RESET} $fp"
        return 0
    fi

    # Bash command
    if [[ "$line" == *'"name":"Bash"'* && "$line" == *'"command"'* ]]; then
        local cmd=""
        if [[ "$line" =~ '"command":"([^"]*)"' ]]; then
            cmd="${match[1]:0:80}"
        fi
        [[ -n "$cmd" ]] && echo -e "${color}[${agent_id}]${RESET} ${YELLOW}⚡ Running${RESET} $cmd"
        return 0
    fi

    # Build success
    if [[ "$line" == *"Tool ran without output or errors"* ]]; then
        echo -e "${color}[${agent_id}]${RESET} ${GREEN}✓ Command succeeded${RESET}"
        return 0
    fi

    # Error in output
    if [[ "$line" == *'"type":"tool_result"'* ]] && [[ "$line" == *error* || "$line" == *Error* || "$line" == *fail* ]]; then
        echo -e "${color}[${agent_id}]${RESET} ${RED}✗ Error detected${RESET} (check output file for details)"
        return 0
    fi

    # Agent completed
    if [[ "$line" == *'"type":"result"'* ]]; then
        echo -e "${color}[${agent_id}]${RESET} ${GREEN}★ AGENT COMPLETED ★${RESET}"
        return 0
    fi

    return 1
}

# Main loop
while true; do
    NOW=$(date +%s)

    for i in {1..${#AGENT_IDS[@]}}; do
        f="${AGENT_FILES[$i]}"
        id="${AGENT_IDS[$i]}"
        color="${AGENT_COLORS[$i]}"

        [[ -f "$f" ]] || continue

        current=$(wc -l < "$f" 2>/dev/null | tr -d ' ')
        last="${LAST_LINES[$i]}"

        if [[ "$current" -gt "$last" ]]; then
            LAST_ACTIVE[$i]=$NOW
            start=$((last + 1))
            tail -n +"$start" "$f" | while IFS= read -r line; do
                process_line "$line" "$id" "$color"
            done
            LAST_LINES[$i]=$current
        else
            idle=$(( NOW - LAST_ACTIVE[$i] ))
            if [[ $idle -gt 120 && $idle -lt 123 ]]; then
                echo -e "${color}[${id}]${RESET} ${RED}⚠ No activity for ${idle}s — possibly hanging${RESET}"
            elif [[ $idle -gt 300 && $(( idle % 60 )) -lt 3 ]]; then
                echo -e "${color}[${id}]${RESET} ${RED}⚠ STALLED for $(( idle / 60 ))min${RESET}"
            fi
        fi
    done

    sleep 2
done
