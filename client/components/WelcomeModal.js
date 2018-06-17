import React from 'react';
import { Modal, Button } from 'react-bootstrap';

class WelcomeModal extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      inputName: props.name,
    };

    this.handleChange = this.handleChange.bind(this);
  }

  componentDidMount() {
    window.addEventListener('keyup', (e) => {
      e.preventDefault();
      if (e.keyCode === 13) document.getElementById('btn').click();
    });
  }

  handleChange(e) {
    this.setState({ inputName: e.target.value });
  }

  render() {
    return (
      <div style={styles.backdrop}>
        <div className="static-modal">
          <Modal.Dialog>
            <Modal.Header>
              <Modal.Title>Welcome to Bumper!</Modal.Title>
            </Modal.Header>
            <Modal.Body>
              <h5 align="left"> Instructions: </h5>
              <p align="left"> Navigate your rocketship around space using your keyboard arrow keys. Bump space junk and other players into the black holes to earn points! But make sure to watch out for the black holes yourself or youll get sucked in!</p>
              <div>
                Enter name:
                <input
                  type="text"
                  value={this.state.inputName}
                  onChange={this.handleChange}
                />
              </div>
            </Modal.Body>
            <Modal.Footer>
              <Button
                bsStyle="primary"
                id="btn"
                onClick={() => this.props.onSubmit(this.state.inputName)}
              >
              Start Bumping!
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
    backgroundColor: '#fff',
    borderRadius: 5,
    height: window.innerHeight / 2,
    width: window.innerWidth / 2,
    zIndex: 10,
  },
  buttonLayout: {
    flexDirection: 'row',
  },
};


export default WelcomeModal;
