// ***************************************************************
//      asrt-data.js
// ***************************************************************

var asrtdata_latest = [];
var asrtdata_all = [];
var alreadyLoaded = {};
var latestDateInMillis = null;

var getKey = function(item) {
  return item.url+"|"+item.timestamp;
};

var isAlreadyLoaded = function(item) {
  var key = getKey(item);
  if (typeof alreadyLoaded[key] === "undefined" || alreadyLoaded[key] === null) {
    return false;
  } else {
    return true;
  }
};

var addDataItems = function(items) {
  if (typeof items === "undefined" || items === null) {
    return;
  }

  for (var i = 0; i < items.length; i++) {
    addDataItem(items[i]);
  }

  updateLatestItems(items);
  updateLatestTime(items);
};

var updateLatestItems = function(items) {
  asrtdata_latest = items;
};

var updateLatestTime = function(items) {
  if (typeof items === "undefined" || items === null) {
    return;
  }

  for (var i = 0; i < items.length; i++) {
    var item = items[i];
    var newMillis = new Date(item.timestamp).getTime();
    if (latestDateInMillis == null || newMillis > latestDateInMillis) {
      latestDateInMillis = newMillis;
    }
  }
};

var addDataItem = function(item) {
    if (!isAlreadyLoaded(item)) {
      var key = getKey(item);
      alreadyLoaded[key] = item;
      asrtdata_all[asrtdata_all.length] = item;
    }
};

var getDataItem = function(index) {
  if (typeof asrtdata_all === "undefined" || asrtdata_all === null || index >= asrtData.length || index < 0) {
    return {};
  }

  return asrtdata_all[index];
};

var getData = function() {
  if (typeof asrtdata_latest === "undefined" || asrtdata_latest === null ) {
    return [];
  }

  return asrtdata_latest.sort(sortAsrtData);
};

var sortAsrtData = function(a,b) {
  return ([a.url,b.url].sort()[0] == a.url) ? -1 : 1;
};


var hourInMillis = 1000*60*60;
var dayInMillis = hourInMillis*24;
var weekInMillis = dayInMillis*7;

var getErrorPercentages = function(differenceMillis) {
  var nowMillis = Date.now();

  var total = 0;
  var errors = 0;

  for (var i = 0; i < asrtdata_all.length; i++) {
    var item = asrtdata_all[i];
    var thenMillis = new Date(item.timestamp).getTime();
    if (nowMillis - thenMillis < differenceMillis) {
      total++;
      if (!item.ok) {
        errors++;
      }
    }
  }

  return errors / total;
};


module.exports.Add = addDataItems;
module.exports.Get = getDataItem;
module.exports.Data = getData;
module.exports.GetLatest = function() { return (latestDateInMillis == null) ? "" : new Date(latestDateInMillis).toLocaleString(); };
module.exports.GetCount = function() { return asrtdata_all.length; };
module.exports.ErrorPercentageInLastHour = function() { return getErrorPercentages(hourInMillis); };
module.exports.ErrorPercentageInLastDay = function() { return getErrorPercentages(dayInMillis); };
module.exports.ErrorPercentageInLastWeek = function() { return getErrorPercentages(weekInMillis); };