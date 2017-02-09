exports.migrate = function(properties) {
	properties['properties']['.properties.org']['value'] = 'system';
	return properties;
};