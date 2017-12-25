import React from 'react';
import PropTypes from 'prop-types';

import Button from '../../layout/Button';
import Icons from '../../layout/Icons';
import * as helpers from '../../../state/addImage/helpers';

export default class ConfigureImage extends React.PureComponent {
    static propTypes = {
        organisations: PropTypes.arrayOf(PropTypes.string).isRequired,
        repositories: PropTypes.object.isRequired,
        config: PropTypes.object.isRequired,
        selectedOrganisation: PropTypes.string.isRequired,
        selectedRepository: PropTypes.string.isRequired,
        onCreateImage: PropTypes.func.isRequired,
        onCancel: PropTypes.func.isRequired,
        onUpdateConfig: PropTypes.func.isRequired
    }

    getOrganisationOptions () {
        return this.props.organisations.map(organisation => 
            <option key={ organisation } value={ organisation }>{ organisation }</option>
        );
    }

    addNewBuildTag = () => {
        const newTags = [].concat(this.props.config.tags);
        newTags.push(helpers.newTagObject());

        const newConfig = Object.assign({}, this.props.config, { tags: newTags });
        this.props.onUpdateConfig(newConfig);
    }

    removeTag = (index) => {
        const newConfig = helpers.removeTagAt(this.props.config, index);
        this.props.onUpdateConfig(newConfig);
    }

    renderTagRow = (tag, index) => {
        const removeTag = () => this.removeTag(index);
        return (
            <tr key={ tag.tag }>
                <td>
                    <select defaultValue={ tag.type }>
                        <option value="branch">Branch</option>
                        <option value="tag">Tag</option>
                    </select>
                </td>
                <td>
                    <input 
                        type="text" 
                        placeholder={ tag.type === 'branch' ? 'e.g. master' : '/.*/ This targets all tags' }
                        defaultValue={ tag.name } />
                </td>
                <td>
                    <input
                        type="text"
                        placeholder="e.g. latest"
                        defaultValue={ tag.tag } />
                </td>
                <td>
                    { index > 0 && <Button icon={ Icons.trash } size="small" onClick={ removeTag } /> }
                </td>
            </tr>
        );
    }

    renderTagRows () {
        return this.props.config.tags.map(this.renderTagRow);
    }

    render () {
        return (
            <div className="step-configure-image">
                <div className="form">
                    <div className="entry entry-half">
                        <label htmlFor="organisation">Github user:</label>
                        <select id="organisation" defaultValue={ this.props.selectedOrganisation }>
                            { this.getOrganisationOptions() }
                        </select>
                    </div>
                    <div className="entry entry-half">
                        <label htmlFor="repository">Repository:</label>
                        <input type="text" id="repository" defaultValue={ this.props.selectedRepository } />
                    </div>
                    <div className="entry">
                        <label htmlFor="dockerfilePath">Dockerfile path:</label>
                        <span className="entry-help">Please enter Dockerfile’s path relative to the repository root</span>
                        <input type="text" id="dockerfilePath" placeholder="e.g. docker/Dockerfile" defaultValue={ this.props.config.dockerfilePath } />

                    </div>
                    <div className="entry">
                        <label>Configure build tags</label>
                        <span className="entry-help">Leave docker tag empty to have the same docker tag with the branch name</span>
                        <table className="table-inline">
                           <thead>
                                <tr>
                                    <th>PUSH TYPE</th>
                                    <th>NAME</th>
                                    <th>DOCKER TAG</th>
                                    <th width="100"></th>
                                </tr>
                            </thead>
                            <tbody>
                                { this.renderTagRows() }
                                <tr className="table-footer">
                                    <td colSpan="4">
                                        <Button size="small" text="Add Build" onClick={ this.addNewBuildTag } />
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                    <div className="form-actions">
                        <Button text="Cancel" size="big" secondary onClick={ this.props.onCancel } />
                        <Button text="Create Image" size="big" onClick={ this.props.onCreateImage } />
                    </div>
                </div>
            </div>
        );
    }
}
