import React from 'react';
import TextField from 'material-ui/lib/TextField';

export default React.createClass({
  getPair: function() {
    return this.props.pair || [];
  },
  render: function() {
    return <div className="Profile-Edit">
      <h1>Profile Edit</h1>

      <p>
        <label>full name</label>
        <br/>
        <TextField
          type="text"
          id="fullname"
          name="fullname"
          hintText=""
          />
      </p>

      <p>
        <label>homepage</label>
        <br/>
        <TextField
          type="text"
          id="homepage"
          name="homepage"
          hintText=""
          />
      </p>

      <p>
        <label>github</label>
        <br/>
        <TextField
          type="text"
          id="github"
          name="github"
          />
      </p>

      <p>
      <label>twitter</label>
      <br/>
      <TextField
        type="text"
        id="twitter"
        name="twitter"
        />
      </p>

    </div>;
  }
});
