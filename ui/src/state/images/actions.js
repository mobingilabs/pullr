export const IMAGES_SHOW_PAGE = 'IMAGES_SHOW_PAGE';
export const IMAGES_ADD_SELECT_PROVIDER = 'IMAGES_ADD_SET_PROVIDER';
export const IMAGES_ADD_ENTRY = 'IMAGES_ADD_ENTRY';
export const IMAGES_ADD_DETAIL = 'IMAGES_ADD_DETAIL';

export function showPage(pageNumber) {
    return {
        type: IMAGES_SHOW_PAGE,
        pageNumber
    };
}

export function newImageSetProvider(provider) {
    return {
        type: IMAGES_ADD_SET_PROVIDER,
        provider
    };
}


export function addEntry(image) {
    return {
        type: IMAGES_ADD_ENTRY,
        image
    };
}

export function addDetail(imageDetail) {
    return {
        type: IMAGES_ADD_DETAIL,
        imageDetail
    };
}
