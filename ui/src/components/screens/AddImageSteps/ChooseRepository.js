import React from 'react';
import PropTypes from 'prop-types';

import Button from '../../layout/Button';
import './ChooseRepository.scss';

export default class ChooseRepository extends React.PureComponent {
    static propTypes = {
        selectedOrganisation: PropTypes.string,
        organisations: PropTypes.arrayOf(PropTypes.string).isRequired,
        repositories: PropTypes.object.isRequired,
        onSelectOrganisation: PropTypes.func.isRequired,
        onSelectRepository: PropTypes.func.isRequired
    }

    renderOrganisation = (organisation) => {
        const classes = this.props.selectedOrganisation === organisation ? 'active' : '';
        const onClick = () => this.props.onSelectOrganisation(organisation);

        return (
            <div className="list-item" key={ organisation }>
                <a className={ classes } onClick={ onClick }>
                    { organisation }
                </a>
            </div>
        );
    } 

    renderRepository = (repository) => {
        return (
            <div className="list-item" key={ repository }>
                <div className="repository-name">{ repository }</div>
                <Button size="small" text="Select" onClick={ () => this.props.onSelectRepository(repository) } />
            </div>
        ) 
    }

    render () {
        return (
            <div className="step-choose-repository">
                <div className="organisation-selection">
                    <div className="list-header">ORGANISATION</div>
                    <div className="list">
                        { this.props.organisations.map(this.renderOrganisation) }
                    </div>
                </div>
                <div className="repository-selection">
                    <div className="list-header">REPOSITORIES</div>
                    <div className="searchbar">
                        <input type="search" placeholder="Search" />
                    </div>
                    <div className="list">
                        { (this.props.repositories[this.props.selectedOrganisation] || []).map(this.renderRepository) }
                    </div>
                </div>
            </div>
        );
    }
}
