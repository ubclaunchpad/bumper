import React from 'react';

class Minimap extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
    };

    // this.handleChange = this.handleChange.bind(this);
  }

  componentDidMount() {
  }

  //   handleChange(e) {
  //   }

  render() {
    return (
      <div style={styles.modal}>
        Minimap
      </div>
    );
  }
}

const styles = {
  modal: {
    display: 'flex',
    flexDirection: 'column',
    justifyContent: 'center',
    position: 'fixed',
    top: 0,
    bottom: 0,
    left: 0,
    right: 0,
    backgroundColor: 'rgba(5,225,255,0.3)',
    borderRadius: 5,
    height: window.innerHeight / 8,
    width: window.innerWidth / 8,
    padding: 50,
    zIndex: 10,
  },
  buttonLayout: {
    flexDirection: 'row',
  },
};


export default Minimap;
