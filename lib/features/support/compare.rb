class CompareResult
  attr_reader :reasons
  attr_writer :equal
  attr_reader :keys

  def initialize
    @equal = true
    @keys = []
    @reasons = []
  end

  def equal?
    @equal
  end

  def keypath
    keys.length > 0 ? keys.reverse.join(".") : "<root>"
  end
end

# Compares two objects for value equality, traversing to compare each
# nested object
def value_compare obj1, obj2, result=nil
  result = CompareResult.new unless result
  return result if obj1 == "IGNORE"
  unless obj1.class == obj2.class
    result.equal = false
    result.reasons << "Object types differ - expected '#{obj1.class}', received '#{obj2.class}'"
    return result
  end

  case obj1
  when Array
    array_compare(obj1, obj2, result)
  when Hash
    hash_compare(obj1, obj2, result)
  when String
    string_compare(obj1, obj2, result)
  else
    unless result.equal = (obj1 == obj2)
      result.reasons << "#{obj1} is not equal to #{obj2}"
    end
  end
  result
end

def array_compare array1, array2, result
  unless array1.length == array2.length
    result.equal = false
    result.reasons << "Expected #{array1.length} items in array, received #{array2.length}"
    return
  end

  array1.each_with_index do |obj1, index|
    value_compare(obj1, array2[index], result)
    unless result.equal?
      result.keys << "#{index}"
      break
    end
  end
end

def hash_compare hash1, hash2, result
  unless hash1.keys.length == hash2.keys.length
    result.equal = false
    missing = hash1.keys - hash2.keys
    unexpected = hash2.keys - hash1.keys
    result.reasons << "Missing keys from hash: #{missing.join(',')}" unless missing.empty?
    result.reasons << "Unexpected keys in hash: #{unexpected.join(',')}" unless unexpected.empty?
    return
  end

  hash1.each do |key, value|
    value_compare(value, hash2[key], result)
    unless result.equal?
      result.keys << key
      break
    end
  end
end

def string_compare template, str2, result
  unless template == str2 or str2 =~ /#{template}/
    result.equal = false
    result.reasons << "'#{str2}' does not match '#{template}'"
  end
end
