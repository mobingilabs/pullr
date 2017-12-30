import * as notificationActions from '../notifications/actions';

export const IMAGES_SHOW_PAGE = 'IMAGES_SHOW_PAGE';
export const IMAGES_ADD_IMAGE = 'IMAGES_ADD_IMAGE';
export const IMAGES_UPDATE_IMAGE = 'IMAGES_UPDATE_IMAGE';

export function showPage(pageNumber) {
    return {
        type: IMAGES_SHOW_PAGE,
        pageNumber
    };
}

export function addImage(image) {
    return {
        type: IMAGES_ADD_IMAGE,
        image
    };
}

export function updateImage(imageName, imageData) {
    return {
        type: IMAGES_UPDATE_IMAGE,
        imageName,
        imageData
    }
}

export function saveChanges(imageName, imageData) {
    return (dispatch) => {
        dispatch(updateImage(imageName, imageData));

        const message = `${imageName} successfully updated.`;
        dispatch(notificationActions.show('images-update-image', message, 'success'));
    };
}
