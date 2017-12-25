import imagesReducer from './images/reducer';
import addImageReducer from './addImage/reducer';

export default function (state, action) {
    return Object.assign({}, state, {
        images: imagesReducer(state.images, action),
        addImage: addImageReducer(state.addImage, action)
    });
}
