// Generated by wp-to-11ty
const { DateTime } = require("luxon");
const slugify = require('@sindresorhus/slugify');
const util = require('util');
const syntaxHighlight = require('eleventy-plugin-highlightjs');

const categories = require('./src/_data/categories.json');
const post_tags = require('./src/_data/post_tags.json');

module.exports = function(eleventy) {
  eleventy.addPassthroughCopy("./src/wp-content");
  eleventy.addPassthroughCopy("./src/static");
  eleventy.addPassthroughCopy("./src/.well-known");
  eleventy.addPassthroughCopy("./src/_redirects");

  eleventy.addPlugin(syntaxHighlight);

  eleventy.addFilter("strftime", function(dateObj, format) {
    return DateTime.fromJSDate(dateObj, {zone: 'utc'}).toFormat(format);
  });

  eleventy.addFilter('console', function(value) {
    return util.inspect(value);
  });

  eleventy.addFilter('wp_tag_slug', function (value) {
    return post_tags[value] || slugify(value);
  });

  eleventy.addFilter('wp_category_slug', function (value) {
    return categories[value]?.nice_name || slugify(value);
  });

  eleventy.addFilter('ref', function (name) {
    return this.getVariables()[name];
  });

  eleventy.addCollection("page", function (collections) {
    return collections.getAllSorted().filter(function (item) {
      return "page" == item.data.type
    });
  });

  eleventy.addCollection("post", function (collections) {
    return collections.getAllSorted().filter(function (item) {
      return "page" != item.data.type
    });
  });

  // a brutish way of getting all of the tags used on any post
  eleventy.addCollection("tags", function (collections) {
    return [...new Set(collections.getAll().map(post => post.data.tags || []).flat())].sort()
  })

  eleventy.addCollection("category", function (collections) {
    const categorized = {};
    collections.getAllSorted().forEach(item => {
      if(item.data.category) {
        if(Array.isArray(item.data.category)) {
          item.data.category.forEach(category => {
            categorized[category] = categorized[category] || [];
            categorized[category].push(item);
          });
        } else {
          categorized[item.data.category] = categorized[item.data.category] || [];
          categorized[item.data.category].push(item);
        }
      }
    });
    return categorized;
  });

  eleventy.addCollection("categories", function (collections) {
    return collections.getAllSorted()
      .map(function (item) { return item.data.category })
      .filter(function (category) { return !!(category); })
      .map(function (item) { return item[0] })
      .filter(function (category, index, arr) { return arr.indexOf(category) == index; });
  });

  // Based on https://github.com/pdehaan/11ty-yearly-archives
  eleventy.addCollection("postsByYear", collection => {
    const data = {};

    collection.getAllSorted()
      .reverse()
      .filter(function (item) {
        return "page" != item.data.type
      })
      .forEach(post => {
        const year = post.date.getFullYear();
        const yearPosts = data[year] || [];
        yearPosts.push(post);
        data[year] = yearPosts;
      });

    return data;
  });

  return {};
};
