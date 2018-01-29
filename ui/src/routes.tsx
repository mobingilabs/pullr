import * as React from 'react';
import { Route } from 'mobx-router';

import ImageScreen from './components/screens/Images';

export default {
    imageDetail: new Route({
        path: '/images/:imageName',
        component: <ImageScreen />
    }),

    images: new Route({
        path: '/images',
        component: <ImageScreen />
    }),

}
