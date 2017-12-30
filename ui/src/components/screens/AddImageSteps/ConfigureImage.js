import React from 'react';
import PropTypes from 'prop-types';

import Button from '../../layout/Button';
import Icons from '../../layout/Icons';
import ImageEditor from '../../widgets/ImageEditor';
import * as helpers from '../../../state/addImage/helpers';

export default class ConfigureImage extends React.Component {
    static propTypes = {
        image: PropTypes.object.isRequired,
        organisations: PropTypes.arrayOf(PropTypes.string).isRequired,
        repositories: PropTypes.arrayOf(PropTypes.string).isRequired,
        onCancel: PropTypes.func.isRequired,
        onCreateImage: PropTypes.func.isRequired,
        onUpdateImage: PropTypes.func.isRequired
    }

    render() {
        return (
            <div className="step-configure-image">
                <ImageEditor image={this.props.image} repositories={this.props.repositories} organisations={this.props.organisations} onUpdateImage={this.props.onUpdateImage}>
                    <div className="form-actions">
                        <Button text="Cancel" size="big" secondary onClick={this.props.onCancel} />
                        <Button text="Create Image" size="big" onClick={this.props.onCreateImage} />
                    </div>
                </ImageEditor>
            </div>
        );
    }
}
