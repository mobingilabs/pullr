export function getById(state, id) {
    return state.notifications[id] || null;
}
