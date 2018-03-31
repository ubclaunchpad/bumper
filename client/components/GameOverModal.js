import React from 'react';
import App from '../App.js';

class GameOverModal extends React.Component {


      render() {
    return (
      <div style={styles.backdrop}>
        <div style={styles.modal} onload="loadText()">
            <b>GAME OVER</b>
            <span style={styles.modalBody}>
            <b> Time alive:</b><br></br><span id ="elapsedTime"></span>
            <b> Points earned: </b><br></br>
             <b> Final ranking: </b>
            </span>
          <button style={styles.restartButton} onClick={this.props.onClose}>
                Restart game
          </button>
        </div>
      </div>
    );
  }
}

function loadText() {
    let time = "1";
    document.getElementById("elapsedTime").innerHTML = time;
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
    fontSize: 20,
    fontFamily: 'Verdana',
  },
  modalBody: {
    textAlign: 'left',
    padding: 10,
  },
  restartButton: {
    height: 30,
    fontSize: 20,
    fontFamily: 'Verdana',
  },
};


export default GameOverModal;
