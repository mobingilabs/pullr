import * as actions from './actions';

export default function (state, action) {
    switch (action.type) {
        case actions.NOTIFICATIONS_SHOW:
            return showNotification(state, action);
        case actions.NOTIFICATIONS_REMOVE:
            return removeNotification(state, action);
        default:
            return state;
    }
}

function removeNotification(state, action) {
    const newState = Object.assign({}, state);
    delete newState[action.id];
    return newState;
}

function showNotification(state, action) {
    return Object.assign({}, state, {
        [action.id]: { message: action.message, type: action.notificationType }
    });
}
