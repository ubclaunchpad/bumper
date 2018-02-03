import React from 'react';

export default class Junk extends React.Component {
  constructor(mass, position, velocity, pointVal, lastBumped) {
    this.mass = mass;
	this.position = position;
	this.velocity = velocity;
	this.pointVal = pointVal;
	this.lastBumped = lastBumped;
  }

  move() {
	  
  }
  
}