export const NOTIFICATIONS_SHOW = 'NOTIFICATION_SHOW';
export const NOTIFICATIONS_REMOVE = 'NOTIFICATIONS_REMOVE';

export function show(id, message, type = 'information') {
    return {
        type: NOTIFICATIONS_SHOW,
        id,
        message,
        notificationType: type
    };
}

export function remove(id) {
    return {
        type: NOTIFICATIONS_REMOVE,
        id
    };
}
