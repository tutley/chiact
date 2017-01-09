var React = require('react');
var ReactRouter = require('react-router');
var Link = ReactRouter.Link;
var MainContainer = require('../containers/MainContainer');

var Home = React.createClass({
  render: function () {
    return (
      <MainContainer>
        <h1>chiact</h1>
        <p className='lead'>Some fancy Motto</p>

      </MainContainer>
    )
  }
});

module.exports = Home;
