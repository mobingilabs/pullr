export function defaultBuildObject() {
    return {
        type: 'branch', name: 'master', tag: 'latest'
    };
}

export function newBuildObject() {
    return {
        type: 'branch', name: '', tag: ''
    };
}
