#!/bin/bash
# HostManager Bash补全脚本
# 使用方法: source completion.bash

_hostmanager_completion() {
    local cur prev commands
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    # 主要命令列表
    commands="connect c list ls l status s search history h favorites fav f groups g init add-host help version"
    
    case "${prev}" in
        hostmanager|hm)
            # 第一级命令补全
            COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
            return 0
            ;;
        connect|c)
            # 连接命令：补全主机名和IP
            local hosts=$(hostmanager list 2>/dev/null | grep -E '^\s+' | sed 's/.*(\([^@]*\)@\([^:]*\):.*/\1 \2/' | tr '\n' ' ')
            COMPREPLY=( $(compgen -W "${hosts}" -- ${cur}) )
            return 0
            ;;
        status|s)
            # 状态命令：补全主机名和IP  
            local hosts=$(hostmanager list 2>/dev/null | grep -E '^\s+' | sed 's/.*(\([^@]*\)@\([^:]*\):.*/\1 \2/' | tr '\n' ' ')
            COMPREPLY=( $(compgen -W "${hosts}" -- ${cur}) )
            return 0
            ;;
        list|ls|l)
            # 列表命令选项
            COMPREPLY=( $(compgen -W "--groups --favorites -g -f" -- ${cur}) )
            return 0
            ;;
        *)
            # 默认情况：如果在连接相关命令后，提供主机补全
            if [[ ${COMP_WORDS[1]} == "connect" ]] || [[ ${COMP_WORDS[1]} == "c" ]] || [[ ${COMP_WORDS[1]} == "status" ]] || [[ ${COMP_WORDS[1]} == "s" ]]; then
                local hosts=$(hostmanager list 2>/dev/null | grep -E '^\s+' | sed 's/.*(\([^@]*\)@\([^:]*\):.*/\1 \2/' | tr '\n' ' ')
                COMPREPLY=( $(compgen -W "${hosts}" -- ${cur}) )
            else
                COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
            fi
            return 0
            ;;
    esac
}

# 注册补全函数
complete -F _hostmanager_completion hostmanager
complete -F _hostmanager_completion hm

# 为便捷别名也添加补全
complete -F _hostmanager_completion hmc
complete -F _hostmanager_completion hms