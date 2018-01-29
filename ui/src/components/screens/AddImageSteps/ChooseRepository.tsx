import * as React from 'react';
import * as PropTypes from 'prop-types';
import { observable, action } from 'mobx';
import { observer, inject } from 'mobx-react';

import Button from '../../layout/Button';
import Image from '../../../state/models/Image';
import WizardStep from '../../widgets/wizard/WizardStep';
import SourceApi from '../../../libs/api/SourceApi';

import './ChooseRepository.scss';

interface Props {
    image: Image;
    next: Function;
    sourceApi?: SourceApi;
}

@inject('sourceApi')
@observer
export default class ChooseRepository extends React.Component<Props> {
    @observable sourceOwners: Array<string> = [];
    @observable sourceRepositories: Array<string> = [];

    componentWillMount() {
        this.props.sourceApi.getOwners(this.props.image.sourceProvider)
            .then(this.setSourceOwners);
    }

    selectSourceOwner = (sourceOwner: string) => {
        this.props.image.sourceOwner = sourceOwner;
        this.props.sourceApi.getRepositories(this.props.image.sourceProvider, this.props.image.sourceOwner)
            .then(this.setSourceRepositories);
    }

    selectSourceRepository = (repository: string) => {
        this.props.image.sourceRepository = repository;
        this.props.next();
    }

    @action.bound
    setSourceOwners(owners: Array<string>) {
        this.sourceOwners = owners;
    }

    @action.bound
    setSourceRepositories(repositories: Array<string>) {
        this.sourceRepositories = repositories;
    }

    renderSourceOwner = (sourceOwner: string) => {
        const classes = this.props.image.sourceOwner === sourceOwner ? 'active' : '';
        return (
            <div className="list-item" key={ sourceOwner }>
                <a className={ classes } onClick={ () => this.selectSourceOwner(sourceOwner) }>
                    { sourceOwner }
                </a>
            </div>
        );
    } 

    renderSourceRepository = (sourceRepository: string) => {
        return (
            <div className="list-item" key={ sourceRepository }>
                <div className="repository-name">{ sourceRepository }</div>
                <Button size="small" text="Select" onClick={ () => this.selectSourceRepository(sourceRepository) } />
            </div>
        ) 
    }

    render () {
        return (
            <WizardStep title="Choose a repository from Github">
                <div className="step-choose-repository">
                    <div className="organisation-selection">
                        <div className="list-header">ORGANISATION</div>
                        <div className="list">
                            { this.sourceOwners.map(this.renderSourceOwner) }
                        </div>
                    </div>
                    <div className="repository-selection">
                        <div className="list-header">REPOSITORIES</div>
                        <div className="searchbar">
                            <input type="search" placeholder="Search" />
                        </div>
                        <div className="list">
                            { this.sourceRepositories.map(this.renderSourceRepository) }
                        </div>
                    </div>
                </div>
            </WizardStep>
        );
    }
}
