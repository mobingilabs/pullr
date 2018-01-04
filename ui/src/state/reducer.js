import imagesReducer from './images/reducer';
import addImageReducer from './addImage/reducer';
import notificationsReducer from './notifications/reducer';
import historyReducer from './history/reducer';

export default function (state, action) {
    return Object.assign({}, state, {
        images: imagesReducer(state.images, action),
        addImage: addImageReducer(state.addImage, action),
        notifications: notificationsReducer(state.notifications, action),
        history: historyReducer(state.history, action)
    });
}
