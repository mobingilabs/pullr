import * as React from 'react';
import { observer } from 'mobx-react';
import AsyncCmd from '../../libs/asyncCmd';

interface Props<T> {
    cmd: AsyncCmd<T>;
}

@observer
export default class InProgress<T> extends React.Component<Props<T>> {
    render() {
        if (this.props.cmd.inProgress) {
            return this.props.children;
        }

        return null;
    }
}

