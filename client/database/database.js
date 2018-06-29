import firebase from 'firebase/app';
import 'firebase/database';

const config = {
  apiKey: 'AIzaSyA4CbYttLND1GH-uoLF523KYkn4tadF6rY',
  authDomain: 'bumperdb-d7f48.firebaseapp.com',
  databaseURL: 'https://bumperdb-d7f48.firebaseio.com',
  projectId: 'bumperdb-d7f48',
  storageBucket: 'bumperdb-d7f48.appspot.com',
  messagingSenderId: '234111044340',
};

// When setting up a second project for Dev vs Prod:
// const config = process.env.NODE_ENV === 'production'
//   ? prodConfig
//   : devConfig;

if (!firebase.apps.length) {
  firebase.initializeApp(config);
}

const db = firebase.database();

export function registerNewTesterEvent() {
  db.ref('leaderboard/Testers/').on('child_added', (snapshot) => {
    console.log(snapshot);
  });
}

export function registerTesterUpdateEvent() {
  db.ref('leaderboard/Testers/').on('child_changed', (snapshot) => {
    console.log(snapshot);
  });
}

export function getDataOnce() {
  db.ref('leaderboard/Testers/').once('value').then((snapshot) => {
    console.log(snapshot);
  });
}
