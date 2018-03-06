import * as React from 'react';
import Image, { Statuses, Causes } from '../../state/models/Image';
import { observer } from 'mobx-react';

interface Props {
    image: Image;
}

@observer
export default class ImageStatus extends React.Component<Props> {
    render() {
        if (!this.props.image.status) {
            return null;
        }

        let classes = ['status', 'img-status']
        let text = '';
        switch (this.props.image.status.cause) {
            case Causes.BuildSuccess:
                classes.push('success');
                text = 'Build Succeed';
                break;

            case Causes.BuildFail:
                classes.push('fail');
                text = 'Build Failed';
                break;

            case Causes.Delete:
                text = 'Image Deleted';

            default:
                text = 'Ready';
        }

        if (this.props.image.status.name === Statuses.Building) {
            text = 'Building';
        }

        let classNames = classes.join(' ');
        return <div className={classNames}>{text}</div>;
    }
}
