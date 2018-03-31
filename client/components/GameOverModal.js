import React from 'react';
import App from '../App.js';

class GameOverModal extends React.Component {
      render() {
    return (
      <div style={styles.backdrop}>
        <div style={styles.modal}>
            <b>Elapsed time in game:</b>  <b>Points earned: </b> <b>Final ranking: </b>
            
          <button style={styles.restartButton} onClick={this.props.onClose}>
                Restart game
          </button>
        </div>
      </div>
        //document.getElementById("time").innerHTML = ELAPSED_TIME;
    );
  }
}

const styles = {
  backdrop: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    position: 'fixed',
    top: 0,
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: 'rgba(5,225,255,0.3)',
    padding: 50,
    zIndex: 10,
  },
  modal: {
    display: 'flex',
    flexDirection: 'column',
    alignSelf: 'center',
    justifyContent: 'space-evenly',
    backgroundColor: '#fff',
    borderRadius: 5,
    height: window.innerHeight / 3,
    width: window.innerWidth / 2,
    zIndex: 10,
    fontSize: 25,
  },
  restartButton: {
    height: 30,
    fontSize: 20,
  },
};


export default GameOverModal;
