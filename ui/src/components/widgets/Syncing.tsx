import * as React from 'react';
import { Icon } from 'react-fa';
import Icons from '../layout/Icons';

export default class Syncing {
    render() {
        return (
            <div className="syncing">
                <Icon name={Icons.loadingSpinner} /> Sync in progress...
            </div>
        )
    }
}
