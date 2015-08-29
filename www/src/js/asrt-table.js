// ***************************************************************
//      asrt-table.js
// ***************************************************************

var React = require("react");
var $ = require("jquery");
var asrtdata = require("./asrt-data");

var AsrtTableHeader = React.createClass({
  render: function() {
    return (
      <thead>
	      <tr>
		      <th>Result</th>
		      <th>Expected HTTP Status Code</th>
		      <th>URL</th>
		      <th>Timestamp</th>
	      </tr>
      </thead>
    );
  }
});

var AsrtTableBody = React.createClass({
  render: function() {
  	var items = this.props.data.map(function (item) {
  	  var resultRowClass = item.ok ? "" : "danger";
  	  var resultClass = item.ok ? "pass" : "fail";
  	  var resultText = item.ok ? "ok" : "fail";
  	  var timestamp = (typeof item.timestamp !== "undefined" && item.timestamp !== null) ? new Date(item.timestamp).toTimeString() : ""
      return (
        <tr className={resultRowClass}>
	        <td className={resultClass}>{resultText}</td>
	        <td>{item.expectation}</td>
	        <td>{item.url}</td>
	        <td>{timestamp}</td>
	    </tr>
      );
    });

    return (
      <tbody>
      	{items}
      </tbody>
    );
  }
});

var AsrtTable = React.createClass({
  loadDataFromServer: function() {  
    $.ajax({
      url: this.props.url,
      dataType: 'json',
      cache: false,
      success: function(data) {
      	asrtdata.Add(data);
        this.setState({data: asrtdata.Data()});
      }.bind(this),
      error: function(xhr, status, err) {
        //console.error(this.props.url, status, err.toString());
      }.bind(this)
    });
  
  },
  getInitialState: function() {
    return {data: []};
  },
  componentDidMount: function() {
    this.loadDataFromServer();
    setInterval(this.loadDataFromServer, this.props.pollInterval);
  },
  componentDidUpdate: function() {
  	updateDataDependencies();
  },
  render: function() {
    return (
        <table className="table table-striped">
	        <AsrtTableHeader  />
	        <AsrtTableBody data={this.state.data} />
        </table>
    );
  }
});


var renderedAsrtTable = null;
var renderedAsrtStatusCount = null;
var renderedAsrtStatusTimestamp = null;

var initialize = function(url, interval) {
	if (typeof url === "undefined" || url === null || url === "") {
		throw new Error("Cannot initialize: Missing URL")
	}
	if (typeof interval === "undefined" || interval === null || interval === 0) {
		interval = 10000;
	}

	renderedAsrtTable = React.render(
	  <AsrtTable url={url} pollInterval={interval} />,
	  document.getElementById('asrt-table')
	);

	renderedAsrtStatusCount = React.render(
	  <AsrtStatusCount />,
	  document.getElementById('asrt-status-count')
	);

	renderedAsrtStatusTimestamp = React.render(
	  <AsrtStatusTimestamp />,
	  document.getElementById('asrt-status-time')
	);
};


var updateDataDependencies = function() {
	updateLastUpdated();
	updateCount();
};

var updateLastUpdated = function() {
	if (renderedAsrtStatusTimestamp !== null) {
		renderedAsrtStatusTimestamp.setState({data: asrtdata.GetLatest()})
	}
};

var updateCount = function() {
	if (renderedAsrtStatusCount !== null) {
		renderedAsrtStatusCount.setState({data: asrtdata.GetCount()})
	}
};

var AsrtStatusTimestamp = React.createClass({
  getInitialState: function() {
    return {data: ""};
  },
  render: function() {
    return (
      <span>
      	{this.state.data}
      </span>
    );
  }
});

var AsrtStatusCount = React.createClass({
  getInitialState: function() {
    return {data: ""};
  },
  render: function() {
    return (
      <span>
      	{this.state.data}
      </span>
    );
  }
});

module.exports.Initialize = initialize;

