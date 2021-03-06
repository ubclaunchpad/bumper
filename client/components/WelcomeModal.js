import React from 'react';
import { Modal, Button, FormControl, FormGroup, ControlLabel } from 'react-bootstrap';
import countries from '../data/countries.json';

const listOfCountries = [];

class WelcomeModal extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      inputName: props.name,
      country: props.country,
    };

    this.handleChange = this.handleChange.bind(this);
    this.handleSelect = this.handleSelect.bind(this);
  }

  componentWillMount() {
    if (listOfCountries.length === 0) {
      Object.keys(countries).forEach((key) => {
        const option = (
          <option value={key} key={key}>
            {countries[key]} {key.toUpperCase().replace(/./g, char => String.fromCodePoint(char.charCodeAt(0) + 127397))}
          </option>
        );

        listOfCountries.push(option);
      });
    }
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

  handleSelect(e) {
    this.setState({ country: e.target.value });
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

              <form>
                <FormGroup controlId="formBasicText">
                  <ControlLabel> Welcome,</ControlLabel>
                  <FormControl
                    type="text"
                    value={this.state.inputName}
                    onChange={this.handleChange}
                  />
                </FormGroup>
                <FormGroup controlId="formControlsSelect">
                  <ControlLabel> From</ControlLabel>

                  <FormControl componentClass="select" value={this.state.country} onChange={this.handleSelect}>
                    <option value="US" key="US">
                      {countries.US} {'US'.toUpperCase().replace(/./g, char => String.fromCodePoint(char.charCodeAt(0) + 127397))}
                    </option>
                    {listOfCountries}
                  </FormControl>
                </FormGroup>
              </form>
            </Modal.Body>
            <Modal.Footer>
              <Button
                bsStyle="primary"
                id="btn"
                onClick={() => this.props.onSubmit(this.state.inputName, this.state.country)}
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
