import React from 'react';
import Flag from 'react-world-flags';

export default class Leaderboard extends React.Component {
  render() {
    if (!this.props.players) {
      return <div />;
    }

    return (
      <div className="bg-light" style={styles.container}>
        <h2 className="bg-primary text-white p-2">Leaderboard</h2>
        <table className="table">
          <thead>
            <tr>
              <th>Country</th>
              <th>Name</th>
              <th>Score</th>
            </tr>
          </thead>
          <tbody>
            {
              this.props.players.map(p => (
                <tr>
                  <td><Flag code={p.country} height={20} /></td>
                  <td>{p.name}</td>
                  <td>{p.points}</td>
                </tr>
              ))
            }
          </tbody>
        </table>
      </div>
    );
  }
}

const styles = {
  container: {
    position: 'absolute',
    top: 0,
    right: 0,
  },
};
