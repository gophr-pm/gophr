import React from 'react';
import ReactDOM from 'react-dom';
import Voting from './components/Voting';

const pair = ['Trainspotting', '28 Days Later'];

ReactDOM.render(
  <h1>Gophr Skeleton code</h1>
  <Voting pair={pair} />,
  document.getElementById('app')
);
