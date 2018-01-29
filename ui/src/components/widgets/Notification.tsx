import * as React from 'react';
import { observer } from "mobx-react";
import { Icon } from 'react-fa';

import Icons from '../layout/Icons';

import './Notification.scss';

interface Props {
    id: string;
}

@observer
export default class Notification extends React.Component<Props> {
    close = () => {
    }

    render() {
        const notification: any = null;
        if (!notification) {
            return null;
        }

        return (
            <div className={`notification ${notification.type}`}>
                <button onClick={this.close}><Icon name={Icons.close} /></button>
                {notification.message}
            </div>
        );
    }
}
