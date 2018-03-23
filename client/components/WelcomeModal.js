import React from 'react';

class WelcomeModal extends React.Component {
  render() {
    return (
      <div style={styles.backdrop}>
          <div style={styles.modal}>
            Hello
            <button onClick={this.props.onClose}>
                close me
            </button>
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
}


export default WelcomeModal;