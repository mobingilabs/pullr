import * as React from 'react';
import { Icon } from 'react-fa';

import Icons from '../layout/Icons';

export default class LoadingSpinner extends React.PureComponent {
    render() {
        return <Icon name={Icons.loadingSpinner} className='fa-spin loading-spinner' />;
    }
}
