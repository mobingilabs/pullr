import * as React from 'react';
import { runInAction } from 'mobx';
import { observer } from 'mobx-react';

import Button from '../layout/Button';
import Icons from '../layout/Icons';

import Image from '../../state/models/Image';

interface Props {
    image: Image;
    sourceRepositories: Array<string>;
    sourceOwners: Array<string>;
}

@observer
export default class ImageEditor extends React.Component<Props> {
    removeBuild = (buildIndex: number) => {
        this.props.image.removeBuild(buildIndex);
    }

    addBuild = () => {
        this.props.image.addBuild();
    }

    bindInputValue = (field: string) => (e: any) => {
        (this.props.image as any)[field] = e.target.value;
    }

    bindBuildInputValue = (buildIndex: number, field: string) => (e: any) => {
        (this.props.image.builds[buildIndex] as any)[field] = e.target.value;
    }

    render() {
        return (
            <div className="image-editor">
                <div className="form">
                    <div className="entry entry-half">
                        <label htmlFor="sourceOwner">Source user:</label>
                        <input type="text" id="sourceOwner" value={this.props.image.sourceOwner} onChange={this.bindInputValue('sourceOwner')} />
                    </div>
                    <div className="entry entry-half">
                        <label htmlFor="sourceRepository">Repository:</label>
                        <input type="text" id="sourceRepository" value={this.props.image.sourceRepository} onChange={this.bindInputValue('sourceRepository')} />
                    </div>
                    <div className="entry">
                        <label htmlFor="imageName">Image name:</label>
                        <input type="text" id="imageName" value={this.props.image.name} onChange={this.bindInputValue('name')}/>
                    </div>
                    <div className="entry">
                        <label htmlFor="dockerfilePath">Dockerfile path:</label>
                        <span className="entry-help">Please enter Dockerfile’s path relative to the repository root</span>
                        <input type="text" id="dockerfilePath" placeholder="e.g. docker/Dockerfile" value={this.props.image.dockerfilePath} onChange={this.bindInputValue('dockerfilePath')} />
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
                                    <th></th>
                                </tr>
                            </thead>
                            <tbody>
                                {this.props.image.builds.map((build, buildIndex) => 
                                    <tr key={buildIndex}>
                                        <td>
                                            <select value={build.type} onChange={this.bindBuildInputValue(buildIndex, 'type')}>
                                                <option value="branch">Branch</option>
                                                <option value="tag">Tag</option>
                                            </select>
                                        </td>
                                        <td>
                                            <input
                                                type="text"
                                                placeholder={build.type === 'branch' ? 'e.g. master' : '/.*/ This targets all tags'}
                                                value={build.name}
                                                onChange={this.bindBuildInputValue(buildIndex, 'name')} />
                                        </td>
                                        <td>
                                            <input
                                                type="text"
                                                placeholder="e.g. latest"
                                                defaultValue={build.tag}
                                                onChange={this.bindBuildInputValue(buildIndex, 'tag')} />
                                        </td>
                                        <td>
                                            {buildIndex > 0 && <Button icon={Icons.trash} size="small" onClick={ () => this.removeBuild(buildIndex) } />}
                                        </td>
                                    </tr>
                                )}
                                <tr className="table-footer">
                                    <td colSpan={4}>
                                        <Button size="small" text="Add Build" onClick={ this.addBuild } />
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
