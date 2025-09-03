#compdef hostmanager hm
# HostManager Zsh补全脚本

_hostmanager() {
    local context state line
    typeset -A opt_args

    _arguments -C \
        '1: :->commands' \
        '*: :->args' \
        && return 0

    case "$state" in
        commands)
            local commands; commands=(
                'connect:连接到指定主机'
                'c:连接到指定主机(简写)'
                'list:显示主机列表'
                'ls:显示主机列表(简写)'
                'l:显示主机列表(简写)'
                'status:检查主机状态'
                's:检查主机状态(简写)'
                'search:搜索主机'
                'history:显示连接历史'
                'h:显示连接历史(简写)'
                'favorites:显示收藏夹'
                'fav:显示收藏夹(简写)'
                'f:显示收藏夹(简写)'
                'groups:按分组显示'
                'g:按分组显示(简写)'
                'init:生成配置文件模板'
                'add-host:交互式添加新主机'
                'help:显示帮助信息'
                'version:显示版本信息'
            )
            _describe 'commands' commands
            ;;
        args)
            case "${words[2]}" in
                connect|c|status|s)
                    # 获取主机名列表进行补全
                    local hosts; hosts=($(hostmanager list 2>/dev/null | grep -E '^\s+' | sed 's/.*(\([^@]*\)@\([^:]*\):.*/\1 \2/' | tr '\n' ' '))
                    _describe 'hosts' hosts
                    ;;
                list|ls|l)
                    local options; options=(
                        '--groups:按分组显示'
                        '-g:按分组显示(简写)'
                        '--favorites:仅显示收藏'
                        '-f:仅显示收藏(简写)'
                    )
                    _describe 'options' options
                    ;;
                search)
                    _message '搜索关键词'
                    ;;
            esac
            ;;
    esac
    
    return 1
}

# 为便捷别名也添加补全
compdef _hostmanager hmc
compdef _hostmanager hms
compdef _hostmanager hml
compdef _hostmanager hmg
compdef _hostmanager hmf

_hostmanager "$@"