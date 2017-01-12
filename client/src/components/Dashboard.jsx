import React, { PropTypes } from 'react';
import { Card, CardTitle, CardText } from 'material-ui/Card';


const Dashboard = ({ name }) => (
  <Card className="container">
    <CardTitle
      title="Dashboard"
      subtitle="You should get access to this page only after authentication."
    />

  {name && <CardText style={{ fontSize: '16px', color: 'green' }}>Hello, {name}, you have successfully used the API with JWT authentication.</CardText>}
  </Card>
);

Dashboard.propTypes = {
  name: PropTypes.string.isRequired
};

export default Dashboard;
