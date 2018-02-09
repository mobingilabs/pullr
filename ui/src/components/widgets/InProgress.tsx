import * as React from 'react';
import { observer } from 'mobx-react';
import AsyncCmd from '../../libs/asyncCmd';

interface Props {
    cmd: AsyncCmd<any, any, any, any, any, any>;
}

@observer
export default class InProgress<T> extends React.Component<Props> {
    render() {
        if (this.props.cmd.inProgress) {
            return this.props.children;
        }

        return null;
    }
}

