import React from 'react';
import PropTypes from 'prop-types';

import Button from '../layout/Button';
import Icons from '../layout/Icons';
import * as helpers from '../../state/addImage/helpers';

export default class ImageEditor extends React.Component {
    static propTypes = {
        image: PropTypes.object.isRequired,
        onUpdateImage: PropTypes.func.isRequired,
        editOnly: PropTypes.bool
    }

    handleOrganisationChange = (e) => {
        const newImage = Object.assign({}, this.props.image);
        newImage.organisation = e.target.value;
        this.props.onUpdateImage(newImage);
    }

    handleRepositoryChange = (e) => {
        const newImage = Object.assign({}, this.props.image);
        newImage.repository = e.target.value;
        this.props.onUpdateImage(newImage);
    }

    handleNameChange = (e) => {
        const newImage = Object.assign({}, this.props.image);
        newImage.name = e.target.value;
        this.props.onUpdateImage(newImage);
    }

    handleDockerfilePathChange = (e) => {
        const newImage = Object.assign({}, this.props.image);
        newImage.dockerfilePath = e.target.value;
        this.props.onUpdateImage(newImage);
    }

    handleAddBuild = () => {
        const newBuilds = [].concat(this.props.image.builds);
        newBuilds.push(helpers.newBuildObject());

        const newImage = Object.assign({}, this.props.image, { builds: newBuilds });
        this.props.onUpdateImage(newImage);
    }

    handleRemoveBuild = (index) => {
        const newBuilds = [].concat(this.props.image.builds);
        newBuilds.splice(index, 1);
    
        const newImage = Object.assign({}, this.props.image, { builds: newBuilds });
        this.props.onUpdateImage(newImage);
    }

    handleUpdateBuildType = (index, e) => {
        const newBuilds = [].concat(this.props.image.builds);
        newBuilds[index] = Object.assign({}, newBuilds[index]);
        newBuilds[index].type = e.target.value;
        
        const newImage = Object.assign({}, this.props.image, { builds: newBuilds });
        this.props.onUpdateImage(newImage);
    }

    handleUpdateBuildName = (index, e) => {
        const newBuilds = [].concat(this.props.image.builds);
        newBuilds[index] = Object.assign({}, newBuilds[index]);
        newBuilds[index].name = e.target.value;

        const newImage = Object.assign({}, this.props.image, { builds: newBuilds });
        this.props.onUpdateImage(newImage);
    }

    handleUpdateBuildTag = (index, e) => {
        const newBuilds = [].concat(this.props.image.builds);
        newBuilds[index] = Object.assign({}, newBuilds[index]);
        newBuilds[index].tag = e.target.value;

        const newImage = Object.assign({}, this.props.image, { builds: newBuilds });
        this.props.onUpdateImage(newImage);
    }

    render() {
        const { image } = this.props;
        return (
            <div className="image-editor">
                <div className="form">
                    <div className="entry entry-half">
                        <label htmlFor="organisation">Source user:</label>
                        <input type="text" id="organisation" defaultValue={image.organisation} onChange={this.handleOrganisationChange}/>
                    </div>
                    <div className="entry entry-half">
                        <label htmlFor="repository">Repository:</label>
                        <input type="text" id="repository" defaultValue={image.repository} onChange={this.handleRepositoryChange}/>
                    </div>
                    <div className="entry">
                        <label htmlFor="imageName">Image name:</label>
                        <input type="text" id="imageName" defaultValue={image.name} onChange={this.handleNameChange} />
                    </div>
                    <div className="entry">
                        <label htmlFor="dockerfilePath">Dockerfile path:</label>
                        <span className="entry-help">Please enter Dockerfile’s path relative to the repository root</span>
                        <input type="text" id="dockerfilePath" placeholder="e.g. docker/Dockerfile" defaultValue={image.dockerfilePath} onChange={this.handleDockerfilePathChange} />
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
                                {image.builds.map((build, buildIndex) => 
                                    <tr key={buildIndex}>
                                        <td>
                                            <select defaultValue={build.type} onChange={this.handleUpdateBuildType.bind(null, buildIndex)}>
                                                <option value="branch">Branch</option>
                                                <option value="tag">Tag</option>
                                            </select>
                                        </td>
                                        <td>
                                            <input
                                                type="text"
                                                placeholder={build.type === 'branch' ? 'e.g. master' : '/.*/ This targets all tags'}
                                                defaultValue={build.name}
                                                onChange={this.handleUpdateBuildName.bind(null, buildIndex)} />
                                        </td>
                                        <td>
                                            <input
                                                type="text"
                                                placeholder="e.g. latest"
                                                defaultValue={build.tag}
                                                onChange={this.handleUpdateBuildTag.bind(null, buildIndex)} />
                                        </td>
                                        <td>
                                            {buildIndex > 0 && <Button icon={Icons.trash} size="small" onClick={this.handleRemoveBuild.bind(null, buildIndex)} />}
                                        </td>
                                    </tr>
                                )}
                                <tr className="table-footer">
                                    <td colSpan="4">
                                        <Button size="small" text="Add Build" onClick={this.handleAddBuild} />
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                    { this.props.children }
                </div>
            </div>
        );
    }
}
