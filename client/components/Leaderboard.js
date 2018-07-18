import React from 'react';
import {
  Panel,
  ListGroup,
  ListGroupItem,
} from 'react-bootstrap';

import Flag from 'react-world-flags';

export default class Leaderboard extends React.Component {
  render() {
    if (!this.props.players) {
      return <div />;
    }

    return (
      <div style={styles.container}>
        <Panel>
          <Panel.Heading componentClass="h4">
            <Panel.Title>
              <div style={styles.header}>
                Leaderboard
              </div>
            </Panel.Title>
          </Panel.Heading>
          <ListGroup>
            {
              this.props.players.map(p => <ListGroupItem key={p.name}><Flag code="can" height="16" />{p.name} {p.name ? p.points : ''}</ListGroupItem>)
            }
          </ListGroup>
        </Panel>
      </div>
    );
  }
}

const styles = {
  container: {
    position: 'absolute',
    top: 0,
    right: 0,
    backgroundColor: 'white',
  },
  header: {
    paddingTop: 5,
    paddingRight: 10,
    paddingLeft: 10,
  },
};
