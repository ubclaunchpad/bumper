import React from 'react';
import { Modal, Button } from 'react-bootstrap';

// eslint-disable-next-line
class GameOverModal extends React.Component {
  render() {
    const minutes = this.props.finalTime.getMinutes();
    const seconds = this.props.finalTime.getSeconds();
    const timeString = `${minutes}:${seconds < 10 ? `0${seconds}` : seconds}`;
    return (
      <div style={styles.backdrop}>
        <div className="static-modal">
          <Modal.Dialog>
            <Modal.Header>
              <Modal.Title>GAME OVER</Modal.Title>
            </Modal.Header>
            <Modal.Body>
              <div>
                <div><b>Time alive:</b> <span>{timeString}</span></div>
                <div><b>Points earned:</b> <span>{this.props.finalPoints}</span></div>
                <div><b>Final ranking:</b> <span>{this.props.finalRanking}</span></div>
              </div>
            </Modal.Body>
            <Modal.Footer>
              <Button
                bsStyle="primary"
                id="btn"
                onClick={() => this.props.onRestart()}
              >
              Restart Game
              </Button>
            </Modal.Footer>
          </Modal.Dialog>
        </div>
      </div>
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
    justifyContent: 'center',
    alignItems: 'center',
    backgroundColor: '#fff',
    borderRadius: 5,
    height: window.innerHeight / 2,
    width: window.innerWidth / 2,
    zIndex: 10,
    fontSize: 20,
    fontFamily: 'Verdana',
  },
  restartButton: {
    height: 30,
    fontSize: 20,
    fontFamily: 'Verdana',
  },
};

export default GameOverModal;
