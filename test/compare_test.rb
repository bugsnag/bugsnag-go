require 'test_helper'
require_relative '../lib/features/support/compare'

class CompareTest < Test::Unit::TestCase

  def test_different_classes
    result = value_compare({}, [])
    assert_false(result.equal?)
    assert_equal("<root>", result.keypath)
  end

  def test_array_value_types
    result = value_compare(["1"], [1])
    assert_false(result.equal?)
    assert_equal("0", result.keypath)
  end

  def test_array_length
    result = value_compare(["1"], ["1", 1])
    assert_false(result.equal?)
    assert_equal("<root>", result.keypath)
    assert_true(value_compare(["1"], ["1"]).equal?)
  end

  def test_array_nested_types
    assert_false(value_compare(["1", {"color": "grey"}], ["1", {"color": "red"}]).equal?)
    assert_true(value_compare(["1", {"color": "grey"}], ["1", {"color": "grey"}]).equal?)
  end

  def test_nested_object_key_paths
    template = {"breakdown": [1,2,{"animal":"bear"}]}
    actual = {"breakdown": [1,2,{"animal":"fox"}]}
    result = value_compare(template, actual)
    assert_equal("breakdown.2.animal", result.keypath)
    assert_equal(result.reasons, ["'fox' does not match 'bear'"])
    assert_false(result.equal?)
  end

  def test_fuzzy_match_string
    template = {"breakdown": [1,2,{"animal":"^\\w+$"}]}
    actual = {"breakdown": [1,2,{"animal":"fox"}]}
    result = value_compare(template, actual)
    assert_true(result.equal?)

    actual = {"breakdown": [1,2,{"animal":"fox-box"}]}
    result = value_compare(template, actual)
    assert_false(result.equal?)
    assert_equal("breakdown.2.animal", result.keypath)
    assert_equal(result.reasons, ["'fox-box' does not match '^\\w+$'"])
  end

  def test_ignore
    template = {"breakdown": [1,2,{"animal":"IGNORE"}]}
    actual = {"breakdown": [1,2,{"animal":"fox"}]}
    result = value_compare(template, actual)
    assert_true(result.equal?)

    template = {"breakdown": "IGNORE"}
    actual = {"breakdown": [1,2,{"animal":"fox"}]}
    result = value_compare(template, actual)
    assert_true(result.equal?)

    template = "IGNORE"
    actual = {"breakdown": [1,2,{"animal":"fox"}]}
    result = value_compare(template, actual)
    assert_true(result.equal?)
  end
end
