import * as React from 'react';
import * as PropTypes from 'prop-types';
import { observable, action, transaction } from 'mobx';
import { observer, inject } from 'mobx-react';
import { Icon } from 'react-fa';

import SourceStore from '../../../state/SourceStore';

import Button from '../../layout/Button';
import Image from '../../../state/models/Image';
import WizardStep from '../../widgets/wizard/WizardStep';
import Icons from '../../layout/Icons';
import AsyncCmd from '../../../libs/asyncCmd';
import InProgress from '../../widgets/InProgress';
import LoadingSpinner from '../../widgets/LoadingSpinner';

import './ChooseRepository.scss';

interface Props {
    image: Image;
    next: Function;
    sources?: SourceStore;
}

@inject('sources')
@observer
export default class ChooseRepository extends React.Component<Props> {
    searchInput: HTMLInputElement;

    componentWillMount() {
        this.props.sources.selectProvider(this.props.image.repository.provider);
        this.props.sources.loadOrganisations.run().done(this.handleOrganisationsLoaded);
    }

    handleOrganisationsLoaded = () => {
        transaction(() => {
            this.props.image.repository.owner = this.props.sources.currentOrganisations[0];
            this.props.sources.selectOrganisation(this.props.sources.currentOrganisations[0]);
            this.props.sources.loadRepositories.run();
        })
    }

    handleFilterChanged = () => {
        this.props.sources.repositoryFilter = this.searchInput.value;
    }

    selectSourceOwner = (sourceOwner: string) => {
        transaction(() => {
            this.props.image.repository.owner = sourceOwner;
            this.props.sources.selectOrganisation(sourceOwner);
            if (this.props.sources.currentRepositories.length === 0) {
                this.props.sources.loadRepositories.run();
            }
        })
    }

    selectSourceRepository = (repository: string) => {
        this.props.image.repository.name = repository;
        this.props.sources.resetSelections();
        this.props.next();
    }

    renderSourceOwner = (sourceOwner: string) => {
        const classes = this.props.image.repository.owner === sourceOwner ? 'active' : '';
        return (
            <div className="list-item" key={sourceOwner}>
                <a className={classes} onClick={() => this.selectSourceOwner(sourceOwner)}>
                    {sourceOwner}
                </a>
            </div>
        );
    }

    renderSourceRepository = (sourceRepository: string) => {
        return (
            <div className="list-item" key={sourceRepository}>
                <div className="repository-name">{sourceRepository}</div>
                <Button size="small" text="Select" onClick={() => this.selectSourceRepository(sourceRepository)} />
            </div>
        )
    }

    renderLoadingSpinner<T>(cmd: AsyncCmd<T>) {
        return (
            <InProgress cmd={cmd}>
                <div className="list-item loading" key="spinner">
                    <LoadingSpinner />
                </div>
            </InProgress>
        );
    }

    render() {
        return (
            <WizardStep title="Choose a repository from Github">
                <div className="step-choose-repository">
                    <div className="organisation-selection">
                        <div className="list-header">ORGANISATION</div>
                        <div className="list">
                            {this.renderLoadingSpinner(this.props.sources.loadOrganisations)}
                            {this.props.sources.currentOrganisations.map(this.renderSourceOwner)}
                        </div>
                    </div>
                    <div className="repository-selection">
                        <div className="list-header">REPOSITORIES</div>
                        <div className="searchbar">
                            <input type="search" placeholder="Search" value={this.props.sources.repositoryFilter} ref={e => this.searchInput = e} onChange={this.handleFilterChanged} />
                        </div>
                        <div className="list">
                            {this.renderLoadingSpinner(this.props.sources.loadRepositories)}
                            {this.props.sources.currentRepositories.map(this.renderSourceRepository)}
                        </div>
                    </div>
                </div>
            </WizardStep>
        );
    }
}
