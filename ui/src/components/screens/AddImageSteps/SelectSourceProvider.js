import React from 'react';
import PropTypes from 'prop-types';

import Button from '../../layout/Button';
import Card from '../../widgets/Card';

import './SelectSourceProvider.scss';

export default class SelectSourceProvider extends React.PureComponent {
    static propTypes = {
        onSelectProvider: PropTypes.func.isRequired
    }

    selectGithub = () => this.props.onSelectProvider('github');

    render () {
        return (
            <div className="step-select-source-provider">
                <Card icon="github" title="Github" background="#CF2A63" dark>
                    <Button text="Link Account & Use" onClick={ this.selectGithub } />
                </Card>
                <Card icon="bitbucket" title="Bitbucket" disabled>
                    Coming soon...
                </Card> 
                <Card icon="gitlab" title="Gitlab" disabled>
                    Coming soon...
                </Card>
            </div>
        )
    }
}
