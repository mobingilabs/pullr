import * as React from 'react';
import * as PropTypes from 'prop-types';
import { observable, action, transaction } from 'mobx';
import { observer, inject } from 'mobx-react';

import Button from '../../layout/Button';
import Image from '../../../state/models/Image';
import WizardStep from '../../widgets/wizard/WizardStep';
import RootStore from '../../../state/RootStore';

import './ChooseRepository.scss';

interface Props {
    image: Image;
    next: Function;
    store?: RootStore;
}

@inject('store')
@observer
export default class ChooseRepository extends React.Component<Props> {
    searchInput: HTMLInputElement;

    componentWillMount() {
        this.props.store.sources.selectProvider(this.props.image.repository.provider);
        this.props.store.sources.loadOrganisations.run().done(this.handleOrganisationsLoaded);
    }

    handleOrganisationsLoaded = () => {
        transaction(() => {
            this.props.image.repository.owner = this.props.store.sources.currentOrganisations[0];
            this.props.store.sources.selectOrganisation(this.props.store.sources.currentOrganisations[0]);
            this.props.store.sources.loadRepositories.run();
        })
    }

    handleFilterChanged = () => {
        this.props.store.sources.repositoryFilter = this.searchInput.value;
    }

    selectSourceOwner = (sourceOwner: string) => {
        transaction(() => {
            this.props.image.repository.owner = sourceOwner;
            this.props.store.sources.selectOrganisation(sourceOwner);
            if (this.props.store.sources.currentRepositories.length === 0) {
                this.props.store.sources.loadRepositories.run();
            }
        })
    }

    selectSourceRepository = (repository: string) => {
        this.props.image.repository.name = repository;
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

    render() {
        return (
            <WizardStep title="Choose a repository from Github">
                <div className="step-choose-repository">
                    <div className="organisation-selection">
                        <div className="list-header">ORGANISATION</div>
                        <div className="list">
                            {this.props.store.sources.currentOrganisations.map(this.renderSourceOwner)}
                        </div>
                    </div>
                    <div className="repository-selection">
                        <div className="list-header">REPOSITORIES</div>
                        <div className="searchbar">
                            <input type="search" placeholder="Search" value={this.props.store.sources.repositoryFilter} ref={e => this.searchInput = e} onChange={this.handleFilterChanged} />
                        </div>
                        <div className="list">
                            {this.props.store.sources.currentRepositories.map(this.renderSourceRepository)}
                        </div>
                    </div>
                </div>
            </WizardStep>
        );
    }
}
