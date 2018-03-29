import React from 'react';

class WelcomeModal extends React.Component {
  constructor(props) {
    super(props);
    this.state = {inputName: ''};
  }

  
  submitPlayerID

  render() {
    return (
      <div style={styles.backdrop}>
          <div style={styles.modal}>
            Welcome to Bumper
            <div style={styles.buttonLayout}>
              <button onClick={this.props.onClose}>
                  close me
              </button>
              <button onClick={() => {this.setState({inputName: 1})}}>
                  submit
              </button>
            </div>
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
        bottom:0,
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
}


export default WelcomeModal;