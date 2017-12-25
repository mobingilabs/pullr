export function defaultTagObject() {
    return {
        type: 'branch', name: 'master', tag: 'latest'
    };
}

export function newTagObject() {
    return {
        type: 'branch', name: '', tag: ''
    };
}

export function removeTagAt(config, index) {
    let tags = [].concat(config.tags);
    tags.splice(index, 1);

    return Object.assign({}, config, { tags });
}
